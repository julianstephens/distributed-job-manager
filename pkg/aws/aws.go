package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func GetConfig() (*aws.Config, *credentials.StaticCredentialsProvider, error) {
	creds := credentials.NewStaticCredentialsProvider("test", "test", "")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, nil, err
	}
	return &cfg, &creds, nil
}

