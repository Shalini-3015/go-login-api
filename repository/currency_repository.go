package repository

import (
	"errors"

	"go-login-api-task/config"
	"go-login-api-task/models"

	"gorm.io/gorm"
)

type CurrencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository() *CurrencyRepository {
	return &CurrencyRepository{
		db: config.DB,
	}
}

func (r *CurrencyRepository) CreateCurrency(currency *models.Currency) error {
	return r.db.Create(currency).Error
}

func (r *CurrencyRepository) GetActiveCurrency() ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.Where("is_active = ?", true).Find(&currencies).Error
	return currencies, err
}

func (r *CurrencyRepository) GetCurrencyByID(id uint) (*models.Currency, error) {
	var currency models.Currency
	err := r.db.First(&currency, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &currency, err
}

func (r *CurrencyRepository) GetCurrencyByCode(code string) (*models.Currency, error) {
	var currency models.Currency
	err := r.db.Where("code = ?", code).First(&currency).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &currency, err
}

func (r *CurrencyRepository) UpdateCurrency(currency *models.Currency) error {
	return r.db.Save(currency).Error
}

func (r *CurrencyRepository) GetActiveCurrencies() ([]models.Currency, error) {
	var currencies []models.Currency

	err := r.db.
		Where("is_active = ?", true).
		Find(&currencies).Error

	return currencies, err
}
