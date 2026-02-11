package service

import (
	"errors"
	"time"
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"go-login-api-task/models"
	"go-login-api-task/repository"
	"go-login-api-task/dto/external"
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

	return s.repo.CreateExcRate(rate)
}

func (s *ExchangeRateService) GetActiveExchangeRates() ([]models.ExchangeRate, error) {
	return s.repo.GetAllActiveExcRate()
}

func (s *ExchangeRateService) GetExchangeRateByID(id uint) (*models.ExchangeRate, error) {
	rate, err := s.repo.GetExcRateByID(id)
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
	rate *float64,
	isActive *bool,
) error {
	existingRate, err := s.repo.GetExcRateByID(id)
    if err != nil {
        return err
    }
    if existingRate == nil {
        return errors.New("exchange rate not found")
    }

    // Update rate only if provided
    if rate != nil {
        if *rate <= 0 {
            return errors.New("rate must be greater than zero")
        }
        existingRate.Rate = *rate
    }

    // Update active state only if provided
    if isActive != nil {
        existingRate.IsActive = *isActive

        if *isActive {
            existingRate.DeletedAt = nil
        } else {
            now := time.Now()
            existingRate.DeletedAt = &now
        }
    }

    return s.repo.UpdateExchangeRate(existingRate)
}


func (s *ExchangeRateService) DeactivateExchangeRate(id uint) error {
	rate, err := s.repo.GetExcRateByID(id)
	if err != nil {
		return err
	}
	if rate == nil {
		return errors.New("exchange rate not found")
	}
	if !rate.IsActive {
		return errors.New("exchange rate already inactive")
	}

now := time.Now()
	rate.IsActive = false
	rate.DeletedAt = &now
	return s.repo.UpdateExcRate(rate)
}

func (s *ExchangeRateService) FetchAndSyncRates(base string) (map[string]float64, error) {

	base = strings.ToUpper(base)

	
	baseCurrency, err := s.currencyRepo.GetCurrencyByCode(base)
	if err != nil || baseCurrency == nil  {
		return nil, fmt.Errorf("invalid base currency")
	}
	if !baseCurrency.IsActive{
		return nil, fmt.Errorf("base currency is inactive")
	}

	
	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", base)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external rates")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API error")
	}

	
	var apiResponse external.ExchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse external response")
	}

	
	activeCurrencies, err := s.currencyRepo.GetActiveCurrencies()
	if err != nil {
		return nil, err
	}

	currencyMap := make(map[string]uint)
	for _, c := range activeCurrencies {
		currencyMap[strings.ToUpper(c.Code)] = c.ID
	}

	result := make(map[string]float64)

	for code, rateValue := range apiResponse.Rates {

		toCurrencyID, exists := currencyMap[strings.ToUpper(code)]
		if !exists {
			continue
		}

		existingRate, err := s.repo.
			GetExchangeRateByCurrencyIDs(baseCurrency.ID, toCurrencyID)

		if err != nil {
			return nil, err
		}

		if existingRate != nil {
		
			existingRate.Rate = rateValue
			existingRate.IsActive = true
			existingRate.DeletedAt = nil

			if err := s.repo.UpdateExchangeRate(existingRate); err != nil {
				return nil, err
			}
		} else {
			
			newRate := &models.ExchangeRate{
				FromCurrencyID: baseCurrency.ID,
				ToCurrencyID:   toCurrencyID,
				Rate:           rateValue,
				IsActive:       true,
			}

			if err := s.repo.CreateExchangeRate(newRate); err != nil {
				return nil, err
			}
		}

		result[code] = rateValue
		// 🔹 Generate inverse rate (target → base)
if rateValue != 0 {

    inverseRate := 1 / rateValue

    inverseExisting, err := s.repo.
        GetExchangeRateByCurrencyIDs(toCurrencyID, baseCurrency.ID)

    if err != nil {
        return nil, err
    }

    if inverseExisting != nil {

        inverseExisting.Rate = inverseRate
        inverseExisting.IsActive = true
        inverseExisting.DeletedAt = nil

        if err := s.repo.UpdateExchangeRate(inverseExisting); err != nil {
            return nil, err
        }

    } else {

        inverseNew := &models.ExchangeRate{
            FromCurrencyID: toCurrencyID,
            ToCurrencyID:   baseCurrency.ID,
            Rate:           inverseRate,
            IsActive:       true,
        }

        if err := s.repo.CreateExchangeRate(inverseNew); err != nil {
            return nil, err
        }
    }
}

	}

	return result, nil
}
