package controller

import (
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
	"gorm.io/gorm"
)

type Controller struct {
	DB     *gorm.DB
	Config *model.Config
}
