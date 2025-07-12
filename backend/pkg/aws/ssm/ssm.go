package ssm

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/config"
	internalaws "github.com/julianstephens/distributed-task-scheduler/backend/pkg/aws"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
)

var (
	client *awsssm.Client
	err    error
	once   sync.Once
	conf   = config.GetConfig()
)

func GetSSMClient() (*awsssm.Client, error) {
	once.Do(func() {
		client, err = setup()
		if err != nil {
			logger.Fatalf("unable to load application environment: %+v", err)
		}
	})

	return client, err
}

func setup() (*awsssm.Client, error) {
	cfg, _, err := internalaws.GetConfig()
	if err != nil {
		return nil, err
	}

	return awsssm.NewFromConfig(*cfg, func(o *awsssm.Options) {
		o.BaseEndpoint = aws.String(conf.BaseEndpoint)
	}), nil
}
