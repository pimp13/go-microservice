package repository

import (
	"microservice_project/services/user/model"

	"gorm.io/gorm"
)

var _ CommonBehaviorsRepository[model.Model] = (*CommonBehaviorsRepositoryImpl[model.Model])(nil)

type CommonBehaviorsRepository[T model.Model] interface {
	Create(data *T) error
	GetByID(id uint) (*T, error)
	GetAll() ([]T, error)
	Update(id uint, data *T) error
	Delete(id uint) error
}

type CommonBehaviorsRepositoryImpl[T model.Model] struct {
	db *gorm.DB
}

func NewCommonBehaviorsRepository[T model.Model](db *gorm.DB) CommonBehaviorsRepository[T] {
	return &CommonBehaviorsRepositoryImpl[T]{db: db}
}

// Create implements CommonBehaviorsRepository.
func (c *CommonBehaviorsRepositoryImpl[T]) Create(data *T) error {
	return c.db.Create(data).Error
}

// Delete implements CommonBehaviorsRepository.
func (c *CommonBehaviorsRepositoryImpl[T]) Delete(id uint) error {
	var entity T
	return c.db.Delete(&entity, id).Error
}

// GetAll implements CommonBehaviorsRepository.
func (c *CommonBehaviorsRepositoryImpl[T]) GetAll() ([]T, error) {
	var entities []T
	if err := c.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// GetByID implements CommonBehaviorsRepository.
func (c *CommonBehaviorsRepositoryImpl[T]) GetByID(id uint) (*T, error) {
	var entity T
	if err := c.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update implements CommonBehaviorsRepository.
func (c *CommonBehaviorsRepositoryImpl[T]) Update(id uint, data *T) error {
	return c.db.Model(data).Where("id = ?", id).Updates(data).Error
}
