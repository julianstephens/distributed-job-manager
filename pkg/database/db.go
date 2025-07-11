package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/guregu/dynamo/v2"
)

func GetDB() (*dynamo.DB, error) {
	creds := credentials.NewStaticCredentialsProvider("test", "test", "")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, err
	}

	return dynamo.New(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localstack:4566")
	}), nil
}
