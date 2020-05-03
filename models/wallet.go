package models

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	PID       uuid.UUID  `gorm:"column:pid;type:varchar(36);index" json:"pid"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	User      uint       `gorm:"column:user_id;not null;index" json:"user_id"`
	Balance   uint       `gorm:"not null" json:"balance"`
	Asset     string     `gorm:"not null;index" json:"asset"`
}
