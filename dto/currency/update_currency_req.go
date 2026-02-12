package dto
import "time"

type  CurrencyUpdateRequest struct {
	Name     *string `json:"name"`
	Symbol   *string `json:"symbol"`
	IsActive *bool   `json:"is_active"`
	DeletedAt *time.Time `json:"deleted_at"`
}
