package external

type ExchangeRateAPIResponse struct {
	Result   string             `json:"result"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}
