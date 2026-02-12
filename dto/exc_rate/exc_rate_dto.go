

package exchange_rate

type UpdateExchangeRateDTO struct {
	Rate     *float64 `json:"rate"`
	IsActive *bool    `json:"is_active"`
}
