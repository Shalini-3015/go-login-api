package repository

import (
	"errors"

	"go-login-api-task/config"
	"go-login-api-task/models"

	"gorm.io/gorm"
)

type ExchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository() *ExchangeRateRepository {
	return &ExchangeRateRepository{
		db: config.DB,
	}
}

func (r *ExchangeRateRepository) Create(rate *models.ExchangeRate) error {
	return r.db.Create(rate).Error
}

func (r *ExchangeRateRepository) GetAllActive() ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	err := r.db.Where("is_active = ?", true).Find(&rates).Error
	return rates, err
}

func (r *ExchangeRateRepository) GetByID(id uint) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	err := r.db.First(&rate, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rate, err
}

func (r *ExchangeRateRepository) GetActiveRate(fromID, toID uint) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	err := r.db.
		Where("from_currency_id = ? AND to_currency_id = ? AND is_active = ?", fromID, toID, true).
		First(&rate).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rate, err
}

func (r *ExchangeRateRepository) Update(rate *models.ExchangeRate) error {
	return r.db.Save(rate).Error
}

func (r *ExchangeRateRepository) Deactivate(id uint) error {
	return r.db.Model(&models.ExchangeRate{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
