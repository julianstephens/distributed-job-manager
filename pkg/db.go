package pkg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/guregu/dynamo/v2"
)

func GetDB() (*dynamo.DB, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithBaseEndpoint("http://localhost:4566"))
	if err != nil {
		return nil, err
	}

	return dynamo.New(cfg), nil
}
