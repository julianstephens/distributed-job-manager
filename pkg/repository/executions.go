package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	"github.com/scylladb/gocqlx/v3/qb"
)

type ExecutionRepository struct {
	Repository
}

func NewExecutionRepository(db *store.DBSession, logger *graylogger.GrayLogger) *ExecutionRepository {
	return &ExecutionRepository{
		Repository{
			DB:     db,
			Logger: logger,
		},
	}
}

// CreateExecution creates a new job execution in the database.
func (r *ExecutionRepository) CreateExecution(execData models.JobExecution) (jobExecution *models.JobExecution, err error) {
	r.Logger.Info("creating job execution", nil)
	execData.ExecutionID = uuid.New().String()
	execData.Status = models.JobStatusScheduled

	if err = r.DB.Client.Query(models.JobExecutions.Insert()).BindStruct(&execData).ExecRelease(); err != nil {
		r.Logger.Error("unable to create job execution", &err)
		err = errors.New("unable to create job execution")
		return
	}

	var res models.JobExecution
	stmt, names := qb.Select(models.JobExecutions.Name()).Where(qb.Eq("execution_id")).AllowFiltering().ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(execData.ExecutionID).SelectRelease(&res); err != nil {
		r.Logger.Error("unable to get created job execution", &err)
		err = errors.New("unable to get created job execution")
		return
	}

	jobExecution = &res

	return
}

// UpdateExecution updates an existing job execution in the database.
func (r *ExecutionRepository) UpdateExecution(execUpdates models.JobExecutionUpdateRequest, executionId string) (jobExecution *models.JobExecution, err error) {
	r.Logger.Info(fmt.Sprintf("updating job execution %s", executionId), nil)

	var res models.JobExecution
	stmt, names := qb.Select(models.JobExecutions.Name()).Where(qb.Eq("execution_id")).AllowFiltering().ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(executionId).SelectRelease(&res); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to get job execution %s", executionId), &err)
		err = fmt.Errorf("unable to get job execution %s", executionId)
		return
	}

	stmt, names = qb.Delete(models.JobExecutions.Name()).Where(qb.Eq("execution_id")).ToCql()
	if err = r.DB.Client.Query(stmt, names).Bind(executionId).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to delete job execution %s", executionId), &err)
		err = fmt.Errorf("unable to delete job execution %s", executionId)
		return
	}

	err = copier.Copy(&res, &execUpdates)
	if err != nil {
		r.Logger.Error("unable to copy updates to job execution", &err)
		err = errors.New("unable to copy updates to job execution")
		return
	}

	if err = r.DB.Client.Query(models.JobExecutions.Insert()).BindStruct(&res).ExecRelease(); err != nil {
		r.Logger.Error(fmt.Sprintf("unable to recreate job execution %s with updates", executionId), &err)
		err = errors.New("unable to update job execution")
		return
	}

	jobExecution = &res

	return
}
