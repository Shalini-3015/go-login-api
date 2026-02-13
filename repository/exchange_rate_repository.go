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

func (r *ExchangeRateRepository) CreateExcRate(rate *models.ExchangeRate) error {
	return r.db.Create(rate).Error
	
}

func (r *ExchangeRateRepository) GetAllActiveExcRate() ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	err := r.db.Where("is_active = ?", true).Find(&rates).Error
	return rates, err
}

func (r *ExchangeRateRepository) GetExcRateByID(id uint) (*models.ExchangeRate, error) {
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




func (r *ExchangeRateRepository) UpdateExcRate(rate *models.ExchangeRate) error {
	return r.db.Save(rate).Error
}










func (r *ExchangeRateRepository) GetExchangeRateByCurrencyIDs(
	fromID uint,
	toID uint,
) (*models.ExchangeRate, error) {

	var rate models.ExchangeRate

	err := r.db.
		Where("from_currency_id = ? AND to_currency_id = ?", fromID, toID).
		First(&rate).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &rate, err
}

