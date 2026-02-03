package repository

import (
	"go-login-api-task/config"
	"go-login-api-task/models"
	"gorm.io/gorm"
)	
type UserRepository struct {
	DB *gorm.DB
}	
func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: config.DB,
	}
}
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}