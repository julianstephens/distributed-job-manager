package controller

import (
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
)

type Controller struct {
	DB     *store.DBSession
	Config *models.Config
}
