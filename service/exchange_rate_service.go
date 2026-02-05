package service

import (
	"errors"

	"go-login-api-task/models"
	"go-login-api-task/repository"
)

type ExchangeRateService struct {
	repo         *repository.ExchangeRateRepository
	currencyRepo *repository.CurrencyRepository
}

func NewExchangeRateService(
	repo *repository.ExchangeRateRepository,
	currencyRepo *repository.CurrencyRepository,
) *ExchangeRateService {
	return &ExchangeRateService{
		repo:         repo,
		currencyRepo: currencyRepo,
	}
}

func (s *ExchangeRateService) CreateExchangeRate(rate *models.ExchangeRate) error {
	if rate.FromCurrencyID == rate.ToCurrencyID {
		return errors.New("from and to currency must be different")
	}

	if rate.Rate <= 0 {
		return errors.New("exchange rate must be greater than zero")
	}

	fromCurrency, err := s.currencyRepo.GetCurrencyByID(rate.FromCurrencyID)
	if err != nil || fromCurrency == nil || !fromCurrency.IsActive {
		return errors.New("from currency not found or inactive")
	}

	toCurrency, err := s.currencyRepo.GetCurrencyByID(rate.ToCurrencyID)
	if err != nil || toCurrency == nil || !toCurrency.IsActive {
		return errors.New("to currency not found or inactive")
	}

	return s.repo.Create(rate)
}

func (s *ExchangeRateService) GetActiveExchangeRates() ([]models.ExchangeRate, error) {
	return s.repo.GetAllActive()
}

func (s *ExchangeRateService) GetExchangeRateByID(id uint) (*models.ExchangeRate, error) {
	rate, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rate == nil || !rate.IsActive {
		return nil, errors.New("exchange rate not found")
	}
	return rate, nil
}

func (s *ExchangeRateService) UpdateExchangeRate(
	id uint, 
	value float64,
	isActive *bool,
	) error {
	if value <= 0 {
		return errors.New("exchange rate must be greater than zero")
	}

	rate, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if rate == nil {
		return errors.New("exchange rate not found")
	}

	rate.Rate = value

	if isActive != nil {
		rate.IsActive = *isActive
	}

	return s.repo.Update(rate)
}

func (s *ExchangeRateService) DeactivateExchangeRate(id uint) error {
	rate, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if rate == nil {
		return errors.New("exchange rate not found")
	}

	return s.repo.Deactivate(id)
}
