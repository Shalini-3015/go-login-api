package service

import (
	"errors"
	"go-login-api-task/models"
	"go-login-api-task/repository"
	
	"strings"
	"time"
	"context"
	"go-login-api-task/dto/currency"
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
	if currency == nil  {
		return nil, errors.New("currency not found")
	}
	return currency, nil
}

func (s *CurrencyService) UpdateCurrency(
	ctx context.Context,
	id uint,
	req dto.CurrencyUpdateRequest,
) error {

	return s.repo.UpdateCurrency(ctx, id, req)
	}

	
func (s *CurrencyService) DeactivateCurrency(
	ctx context.Context,
	id uint,
) error {

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
	active := false

	req := dto.CurrencyUpdateRequest{
		IsActive:  &active,
		DeletedAt: &now,
	}

	return s.repo.UpdateCurrency(ctx, id, req)
}
