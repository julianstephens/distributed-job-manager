package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
	"github.com/oklog/ulid/v2"
	"github.com/scylladb/gocqlx/v3/qb"
)

type JobRepository struct {
	Repository
}

func NewJobRepository(db *store.DBSession, logger *graylogger.GrayLogger) *JobRepository {
	return &JobRepository{
		Repository{
			DB:     db,
			Logger: logger,
		},
	}
}

// GetJobs retrieves all jobs for a user or all jobs if the user is an admin.
func (r *JobRepository) GetJobs(userId string, isAdmin bool) (jobs *[]models.Job, err error) {
	r.Logger.Debug(fmt.Sprintf("getting jobs for user %s", userId), utils.StringPtr(fmt.Sprintf("isAdmin: %t", isAdmin)))

	q := qb.Select(models.Jobs.Name())

	if !isAdmin {
		q.Where(qb.Eq("user_id")).AllowFiltering()
	}

	stmt, names := q.ToCql()
	transaction := r.DB.Client.Query(stmt, names)

	if !isAdmin {
		transaction.Bind(userId)
	}

	var res []models.Job
	if err = transaction.SelectRelease(&res); err != nil {
		data := map[string]any{
			"userId":  userId,
			"isAdmin": isAdmin,
		}
		r.Logger.ErrorWithData("failed to get jobs", &err, &data)
		return
	}

	jobs = &res

	return
}

// GetJob retrieves a specific job by its ID.
func (r *JobRepository) GetJob(jobId string, userId string, isAdmin bool) (job *models.Job, err error) {
	r.Logger.Debug(fmt.Sprintf("getting job %s for user %s", jobId, userId), utils.StringPtr(fmt.Sprintf("isAdmin: %t", isAdmin)))

	q := qb.Select(models.Jobs.Name())

	if !isAdmin {
		q.Where(qb.Eq("user_id"))
	}

	stmt, names := q.Where(qb.Eq("job_id")).AllowFiltering().ToCql()
	transaction := r.DB.Client.Query(stmt, names)

	if isAdmin {
		transaction.Bind(jobId)
	} else {
		transaction.Bind(userId, jobId)
	}

	var res models.Job
	if err = transaction.Get(&res); err != nil {
		data := map[string]any{
			"jobId":   jobId,
			"userId":  userId,
			"isAdmin": isAdmin,
		}
		r.Logger.ErrorWithData("failed to get job", &err, &data)
		return
	}

	job = &res

	return
}

// CreateJob creates a new job in the database with the provided job data and user ID.
func (r *JobRepository) CreateJob(jobData models.Job, userId string) (job *models.Job, err error) {
	now := time.Now().UTC()

	jobData.CreatedAt = now
	jobData.UpdatedAt = now
	jobData.JobID = ulid.Make().String()
	jobData.UserID = userId
	jobData.RetryCount = 0
	jobData.Status = models.JobStatusPending

	parser := &utils.Parser{}
	if err = parser.Parse(jobData.Payload); err != nil {
		r.Logger.Error("failed to parse job payload", &err)
		err = fmt.Errorf("failed to parse job payload: %w", err)
		return
	}

	supportedLanguages := utils.GetSupportedLanguages()
	for _, block := range parser.Result {
		if supportedLanguages[block.Language] == "" {
			r.Logger.Error(fmt.Sprintf("%s is not a supported code language", block.Language), nil)
			err = fmt.Errorf("%s is not a supported code language", block.Language)
			return
		}
	}
	jobData.Payload = parser.SanitizedInput

	if err = r.DB.Client.Query(models.Jobs.Insert()).BindStruct(jobData).ExecRelease(); err != nil {
		r.Logger.Error("failed to insert job into database", &err)
		err = errors.New("failed to insert job into database")
		return
	}

	job = &jobData

	jobSchedule := models.JobSchedule{
		JobID:       job.JobID,
		NextRunTime: job.ExecutionTime,
	}

	if err = r.DB.Client.Query(models.JobSchedules.Insert()).BindStruct(jobSchedule).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("failed to create job schedule for job %s", job.JobID), &err)
		err = fmt.Errorf("failed to create job schedule")
		return
	}

	return
}

// UpdateJob updates an existing job with new data.
func (r *JobRepository) UpdateJob(jobUpdate models.JobUpdateRequest, jobId string, userId string) (job *models.Job, err error) {
	if jobUpdate.Payload != nil {
		parser := &utils.Parser{}
		if err = parser.Parse(*jobUpdate.Payload); err != nil {
			r.Logger.Error("failed to parse job payload", &err)
			err = errors.New("failed to parse job payload")
			return
		}
		supportedLanguages := utils.GetSupportedLanguages()
		for _, block := range parser.Result {
			if supportedLanguages[block.Language] == "" {
				err = fmt.Errorf("%s is not a supported code language", block.Language)
				return
			}
		}
		jobUpdate.Payload = &parser.SanitizedInput
	}

	r.Logger.Info(fmt.Sprintf("updating job %s for user %s", jobId, userId), nil)

	var res models.Job
	stmt, names := qb.Select(models.Jobs.Name()).Where(qb.Eq("job_id")).AllowFiltering().ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(jobId).Get(&res); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to get job %s", jobId), &err)
		err = errors.New("unable to get job")
		return
	}

	stmt, names = qb.Delete(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id"), qb.Eq("status")).ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(res.JobID, res.UserID, res.Status).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to delete job %s", jobId), &err)
		err = errors.New("unable to delete job")
		return
	}

	err = copier.Copy(&res, &jobUpdate)
	if err != nil {
		r.Logger.Error("unable to get copy updates to job", &err)
		err = errors.New("unable to copy updates to job")
		return
	}

	if err = r.DB.Client.Query(models.Jobs.Insert()).BindStruct(res).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to recreate job %s w/ updates", jobId), &err)
		err = errors.New("unable to update job")
		return
	}

	job = &res

	if jobUpdate.ExecutionTime != nil {
		var jobSchedule models.JobSchedule
		updatedSchedule := models.JobSchedule{
			JobID:       job.JobID,
			NextRunTime: *jobUpdate.ExecutionTime,
		}

		stmt, names = qb.Select(models.JobSchedules.Name()).Where(qb.Eq("job_id")).AllowFiltering().ToCql()
		if err = r.DB.Client.Query(stmt, names).Bind(job.JobID).Get(&jobSchedule); err != nil {
			r.Logger.Error(fmt.Sprintf("unable to get job schedule for job %s", jobId), &err)
			err = errors.New("unable to update job schedule")
			return
		}

		stmt, names = qb.Delete(models.JobSchedules.Name()).Where(qb.Eq("job_id"), qb.Eq("next_run_time")).ToCql()
		if err = r.DB.Client.Query(stmt, names).Bind(job.JobID, jobSchedule.NextRunTime).ExecRelease(); err != nil {
			r.Logger.Error(fmt.Sprintf("unable to delete job schedule for job %s", jobId), &err)
			err = errors.New("unable to update job schedule")
			return
		}

		if err = r.DB.Client.Query(models.JobSchedules.Insert()).BindStruct(updatedSchedule).ExecRelease(); err != nil {
			r.Logger.Error(fmt.Sprintf("unable to recreate job schedule for job %s", job.JobID), &err)
			err = errors.New("unable to update job schedule")
			return
		}
	}
	return
}

// DeleteJob removes a job from the database by its ID.
func (r *JobRepository) DeleteJob(jobId string) (err error) {
	var job models.Job
	stmt, names := qb.Select(models.Jobs.Name()).Where(qb.Eq("job_id")).AllowFiltering().ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(jobId).Get(&job); err != nil {
		r.Logger.Error(fmt.Sprintf("job %s not found", jobId), &err)
		err = fmt.Errorf("job %s not found", jobId)
		return
	}

	r.Logger.Info(fmt.Sprintf("deleting job %s for user %s", jobId, job.UserID), nil)
	stmt, names = qb.Delete(models.Jobs.Name()).Where(qb.Eq("user_id"), qb.Eq("job_id"), qb.Eq("status")).ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(job.UserID, job.JobID, job.Status).ExecRelease(); err != nil {
		r.Logger.ErrorWithData(fmt.Sprintf("unable to delete job %s", jobId), &err, &map[string]any{
			"jobId":  jobId,
			"userId": job.UserID,
			"status": job.Status,
		})
		err = fmt.Errorf("unable to delete job %s", jobId)
		return
	}

	// TODO: Also delete any associated schedules, logs, etc.

	return
}
