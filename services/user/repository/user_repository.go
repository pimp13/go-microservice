package repository

import (
	"microservice_project/services/user/model"

	"gorm.io/gorm"
)

var _ UserRepository = (*UserRepositoryImpl)(nil)

type UserRepository interface {
	CommonBehaviorsRepository[model.User]
	GetByEmail(email string) (*model.User, error)
}

type UserRepositoryImpl struct {
	CommonBehaviorsRepository[model.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		CommonBehaviorsRepository: NewCommonBehaviorsRepository[model.User](db),
		db:                        db,
	}
}

// GetByEmail implements UserRepository.
func (u *UserRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
