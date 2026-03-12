package service

import (
	"errors"

	"gorm.io/gorm"
)

type DefaultCrud[T any, ID comparable] struct {
	db *gorm.DB
}

func NewDefaultCrud[T any, ID comparable](db *gorm.DB) *DefaultCrud[T, ID] {
	return &DefaultCrud[T, ID]{db: db}
}




func (r *DefaultCrud[T, ID]) Create(entity *T) (*T, error) {
	if err := r.db.Create(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *DefaultCrud[T, ID]) Get(id ID) (*T, error) {
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *DefaultCrud[T, ID]) GetAll() ([]T, error) {
	var list []T
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *DefaultCrud[T, ID]) Update(id ID, data *T) (*T, error) {
	var entity T
	if err := r.db.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := r.db.Model(&entity).Updates(data).Error; err != nil {
		return nil, err
	}
	return &entity, nil

}

func (r *DefaultCrud[T, ID]) Delete(id ID) (bool, error) {
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if err := r.db.Delete(&entity).Error; err != nil {
		return false, err
	}
	return true, nil
}
