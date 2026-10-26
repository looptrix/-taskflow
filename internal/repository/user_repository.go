package repository

import (
	"errors"

	"gorm.io/gorm"
	"taskflow/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *UserRepository) IsEmailTaken(email string) (bool, error) {
	var count int64

	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err

}

func (r *UserRepository) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
