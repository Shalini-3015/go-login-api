package models
import (
	
	"time"
)

type LoginUser struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Email     string    `gorm:"column:email;unique;not null" json:"email"`
	Password  string    `gorm:"column:password;not null" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:false" json:"updated_at,omitempty"`
}

func (LoginUser) TableName() string {
	return "login_user"
}
