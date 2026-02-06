package service

import (
	"errors"
	"strings"

	"go-login-api-task/repository"
)

type ConversionResponse struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	Amount          float64 `json:"amount"`
	ExchangeRate    float64 `json:"exchange_rate"`
	ConvertedAmount float64 `json:"converted_amount"`
}

type ConversionService struct {
	currencyRepo *repository.CurrencyRepository
	rateRepo     *repository.ExchangeRateRepository
}

func NewConversionService(
	currencyRepo *repository.CurrencyRepository,
	rateRepo *repository.ExchangeRateRepository,
) *ConversionService {
	return &ConversionService{
		currencyRepo: currencyRepo,
		rateRepo:     rateRepo,
	}
}

func (s *ConversionService) ConvertCurrencyAmt(from, to string, amount float64) (*ConversionResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	fromCurrency, err := s.currencyRepo.GetCurrencyByCode(from)
	if err != nil || fromCurrency == nil  {
		return nil, errors.New("from currency not found")
	}
	if !fromCurrency.IsActive {
		return nil, errors.New("from currency is inactive")
	}


	toCurrency, err := s.currencyRepo.GetCurrencyByCode(to)
	if err != nil || toCurrency == nil  {
		return nil, errors.New("to currency not found ")
	}
	if !toCurrency.IsActive {
		return nil, errors.New("to currency is inactive")
	}

	rate, err := s.rateRepo.GetActiveRate(fromCurrency.ID, toCurrency.ID)
	if err != nil || rate == nil {
		return nil, errors.New("exchange rate not found")
	}

	converted := amount * rate.Rate

	return &ConversionResponse{
		From:            from,
		To:              to,
		Amount:          amount,
		ExchangeRate:    rate.Rate,
		ConvertedAmount: converted,
	}, nil
}
