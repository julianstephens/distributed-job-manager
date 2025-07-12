package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
)

var (
	config *model.Config
	once   sync.Once
	err    error
)

func GetConfig() *model.Config {
	once.Do(func() {
		config, err = setup()
		if err != nil {
			logger.Fatalf("unable to load application environment: %+v", err)
		}
	})

	return config
}

func setup() (*model.Config, error) {
	var conf model.Config

	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
