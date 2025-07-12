package controller

import (
	"github.com/guregu/dynamo/v2"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
)

type Controller struct {
	DB     *dynamo.DB
	Config *model.Config
}
