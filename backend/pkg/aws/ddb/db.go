package ddb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/guregu/dynamo/v2"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/config"
	internalaws "github.com/julianstephens/distributed-task-scheduler/backend/pkg/aws"
)

var conf = config.GetConfig()

func GetDB() (*dynamo.DB, error) {
	cfg, _, err := internalaws.GetConfig()
	if err != nil {
		return nil, err
	}

	return dynamo.New(*cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(conf.BaseEndpoint)
	}), nil
}
