package models

import "time"

type Currency struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(3);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Symbol    string    `gorm:"type:varchar(10)" json:"symbol"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
