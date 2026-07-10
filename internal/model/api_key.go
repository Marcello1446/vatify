package model

import (
	"gorm.io/gorm"
	"time"
)

type APIKey struct {
	gorm.Model
	StripeCustomerID string
	Key              string `gorm:"not null;uniqueIndex"`
	IsActive         bool
	Owner            string `gorm:"not null"`
	Plan             string
	RequestLimit     int
	RequestsUsed     int
	UsageReset       time.Time
}
