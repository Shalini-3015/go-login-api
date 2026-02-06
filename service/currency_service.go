package service

import (
	"errors"
	"strings"

	"go-login-api-task/models"
	"go-login-api-task/repository"
)

type CurrencyService struct {
	repo *repository.CurrencyRepository
}

func NewCurrencyService(repo *repository.CurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (s *CurrencyService) CreateCurrency(currency *models.Currency) error {
	currency.Code = strings.ToUpper(currency.Code)

	if currency.Code == "" || currency.Name == "" {
		return errors.New("currency code and name are required")
	}

	existing, err := s.repo.GetCurrencyByCode(currency.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("currency code already exists")
	}

	return s.repo.CreateCurrency(currency)
}

func (s *CurrencyService) GetAllActiveCurrencies() ([]models.Currency, error) {
	return s.repo.GetActiveCurrency()
}

func (s *CurrencyService) GetCurrencyByID(id uint) (*models.Currency, error) {
	currency, err := s.repo.GetCurrencyByID(id)
	if err != nil {
		return nil, err
	}
	if currency == nil || !currency.IsActive {
		return nil, errors.New("currency not found")
	}
	return currency, nil
}

func (s *CurrencyService) UpdateCurrency(id uint, name, symbol *string, isActive *bool) error {
	currency, err := s.repo.GetCurrencyByID(id)
	if err != nil {
		return err
	}
	if currency == nil {
		return errors.New("currency not found")
	}

	if name != nil {
	currency.Name = *name
}

if symbol != nil {
	currency.Symbol = *symbol
}

if isActive != nil {
	currency.IsActive = *isActive
}


	return s.repo.UpdateCurrency(currency)
}

func (s *CurrencyService) DeactivateCurrency(id uint) error {
	currency, err := s.repo.GetCurrencyByID(id)
	if err != nil {
		return err
	}
	if currency == nil {
		return errors.New("currency not found")
	}

	return s.repo.DeactivateCurrency(id)
}
