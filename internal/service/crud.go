package service

import (
	"context"

	"github.com/guregu/dynamo/v2"
)

func GetAll[T any](db *dynamo.DB, tableName string) (*[]T, error) {
	var result []T

	if err := db.Table(tableName).Scan().All(context.Background(), dynamo.AWSEncoding(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

func FindById[T any](db *dynamo.DB, id string, tableName string) (*T, error) {
	var result T
	if err := db.Table(tableName).Get("id", id).One(context.Background(), dynamo.AWSEncoding(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

func Put[T any](db *dynamo.DB, newResource T, tableName string) (*T, error) {
	if err := db.Table(tableName).Put(dynamo.AWSEncoding(newResource)).Run(context.Background()); err != nil {
		return nil, err
	}

	return &newResource, nil
}

func Delete[T any](db *dynamo.DB, id string, tableName string) error {
	if err := db.Table(tableName).Delete("id", id).Run(context.Background()); err != nil {
		return err
	}

	return nil
}
