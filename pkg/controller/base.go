package controller

import (
	"github.com/julianstephens/distributed-task-scheduler/pkg/models"
	"github.com/julianstephens/distributed-task-scheduler/pkg/store"
)

type Controller struct {
	DB     *store.DBSession
	Config *models.Config
}
