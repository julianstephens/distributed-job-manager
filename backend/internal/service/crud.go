package service

import (
	"gorm.io/gorm"
)

type PreloadOptions struct {
	ShouldPreload bool
	PreloadQuery  string
}

func GetAll[T any](db *gorm.DB) (*[]T, error) {
	var result []T

	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func Find[T any](db *gorm.DB, conditions T, preloadOpts *PreloadOptions) (*T, error) {
	var result T

	if preloadOpts != nil && preloadOpts.ShouldPreload {
		if err := db.Preload(preloadOpts.PreloadQuery).Where(conditions).First(&result).Error; err != nil {
			return nil, err
		}
	}

	if err := db.Where(conditions).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func FindById[T any](db *gorm.DB, id string) (*T, error) {
	var result T
	if err := db.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func Create[T any](db *gorm.DB, newResource T, conditions ...T) (*T, error) {
	if err := db.Create(&newResource).Error; err != nil {
		return nil, err
	}

	return &newResource, nil
}

func Update[T any](db *gorm.DB, id string, updates map[string]any) (T, error) {
	existing, err := FindById[T](db, id)
	if err != nil {
		return *new(T), err
	}

	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return *new(T), err
	}

	return *existing, nil
}

func Delete[T any](db *gorm.DB, id string, res T) (*T, error) {
	var result T
	if err := db.Delete(res, id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
