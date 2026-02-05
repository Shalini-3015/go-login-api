package models

import "time"

type ExchangeRate struct {
	ID uint `gorm:"primaryKey" json:"id"`

	FromCurrencyID uint     `gorm:"not null;index" json:"from_currency_id"`
	FromCurrency   Currency `gorm:"foreignKey:FromCurrencyID;references:ID" json:"-"`

	ToCurrencyID uint     `gorm:"not null;index" json:"to_currency_id"`
	ToCurrency   Currency `gorm:"foreignKey:ToCurrencyID;references:ID" json:"-"`

	Rate float64 `gorm:"type:numeric(12,4);not null" json:"rate"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
