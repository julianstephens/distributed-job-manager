package repository

import (
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
)

type Repository struct {
	DB     *store.DBSession
	Logger *graylogger.GrayLogger
}
