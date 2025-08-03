package controller

import (
	"errors"

	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
)

type Controller struct {
	DB     *store.DBSession
	Config *models.Config
	Logger *graylogger.GrayLogger
}

var ErrMissingJobId = errors.New("no job id provided")
