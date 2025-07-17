package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/julianstephens/distributed-task-scheduler/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/pkg/models"
)

var (
	config *models.Config
	once   sync.Once
	err    error
)

func GetConfig() *models.Config {
	once.Do(func() {
		config, err = setup()
		if err != nil {
			logger.Fatalf("unable to load application environment: %+v", err)
		}
	})

	return config
}

func setup() (*models.Config, error) {
	var conf models.Config

	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
