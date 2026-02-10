package service

import (
	"errors"
	"go-login-api-task/models"
	"go-login-api-task/repository"
	"log"
	"strings"
	"time"
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

func (s *CurrencyService) UpdateCurrency(
	id uint,
	name, symbol *string,
	isActive *bool,
) error {
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
		if *isActive && !currency.IsActive {
			currency.IsActive = true
			currency.DeletedAt = nil
		}
		if !*isActive {
			return errors.New("use delete api to deactivate currency")
		}
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
	if !currency.IsActive {
		return errors.New("currency already inactive")
	}
	now := time.Now()

	currency.IsActive = false
	currency.DeletedAt = &now
	log.Println(">>> setting deleted_at:", currency.DeletedAt)
	return s.repo.UpdateCurrency(currency)
}
