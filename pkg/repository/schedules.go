package repository

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
)

type ScheduleRepository struct {
	Repository
}

func NewScheduleRepository(db *store.DBSession, logger *graylogger.GrayLogger) *ScheduleRepository {
	return &ScheduleRepository{
		Repository{
			DB:     db,
			Logger: logger,
		},
	}
}

// GetSchedules retrieves all job schedules from the database, applying any filters specified in queryParams.
func (r *ScheduleRepository) GetSchedules(queryParams map[string][]string) (jobSchedules *[]models.JobSchedule, err error) {
	r.Logger.Info("retrieving all job schedules", nil)

	var res []models.JobSchedule
	q, err := r.getFilteredQuery(queryParams, models.JobSchedules.Name(), models.JobSchedules.Metadata().Columns)
	if err != nil {
		r.Logger.Error("unable to get filtered query for job schedules", &err)
		err = errors.New("unable to get job schedules")
		return
	}

	if err = q.SelectRelease(&res); err != nil {
		r.Logger.Error("unable to get job schedules from db", &err)
		err = errors.New("unable to get job schedules")
		return
	}

	jobSchedules = &res

	return
}

// GetSchedule retrieves a specific job schedule by its ID from the database.
func (r *ScheduleRepository) GetSchedule(id string) (jobSchedule *models.JobSchedule, err error) {
	r.Logger.Info(fmt.Sprintf("retrieving job schedule %s", id), nil)

	var res models.JobSchedule
	if err = r.DB.Client.Query(models.JobSchedules.Select()).Bind(id).Get(&res); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to get job schedule %s from db", id), &err)
		err = fmt.Errorf("unable to get job schedule %s", id)
		return
	}

	jobSchedule = &res

	return
}

// CreateSchedule inserts a new job schedule into the database.
func (r *ScheduleRepository) CreateSchedule(scheduleData models.JobSchedule) (jobSchedule *models.JobSchedule, err error) {
	r.Logger.Info("creating new job schedule", nil)

	if err = r.DB.Client.Query(models.JobSchedules.Insert()).BindStruct(&scheduleData).ExecRelease(); err != nil {
		r.Logger.Error("unable to create job schedule in db", &err)
		err = errors.New("unable to create job schedule")
		return
	}

	jobSchedule = &scheduleData

	return
}

// UpdateSchedule modifies an existing job schedule in the database based on the provided updates.
func (r *ScheduleRepository) UpdateSchedule(id string, scheduleUpdates models.JobScheduleUpdateRequest) (jobSchedule *models.JobSchedule, err error) {
	r.Logger.Info(fmt.Sprintf("updating job schedule %s", id), nil)

	var existingSchedule models.JobSchedule
	if err = r.DB.Client.Query(models.JobSchedules.Select()).Bind(id).Get(&existingSchedule); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to get job schedule %s in db", id), &err)
		err = fmt.Errorf("unable to get job schedule %s", id)
		return
	}

	if err = r.DB.Client.Query(models.JobSchedules.Delete()).BindMap(map[string]any{"job_id": id, "next_run_time": existingSchedule.NextRunTime}).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to delete job schedule %s in db", id), &err)
		err = fmt.Errorf("unable to update job schedule %s", id)
		return
	}

	if err = copier.Copy(&existingSchedule, &scheduleUpdates); err != nil {
		r.Logger.Error("unable to copy schedule updates", &err)
		err = errors.New("unable to copy schedule updates")
		return
	}

	if err = r.DB.Client.Query(models.JobSchedules.Insert()).BindStruct(&existingSchedule).ExecRelease(); err != nil {
		r.Logger.Error("unable to update job schedule in db", &err)
		err = errors.New("unable to update job schedule")
		return
	}

	jobSchedule = &existingSchedule
	return
}

// DeleteSchedule removes a job schedule from the database by its ID.
func (r *ScheduleRepository) DeleteSchedule(id string) (err error) {
	r.Logger.Info(fmt.Sprintf("deleting job schedule %s", id), nil)

	var existingSchedule models.JobSchedule
	if err = r.DB.Client.Query(models.JobSchedules.Select()).Bind(id).Get(&existingSchedule); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to get job schedule %s in db", id), &err)
		err = fmt.Errorf("unable to get job schedule %s", id)
		return
	}

	if err = r.DB.Client.Query(models.JobSchedules.Delete()).BindMap(map[string]interface{}{"job_id": id, "next_run_time": existingSchedule.NextRunTime}).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to delete job schedule %s in db", id), &err)
		err = fmt.Errorf("unable to delete job schedule %s", id)
		return
	}

	return
}
