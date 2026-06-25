package model

import (
	"gorm.io/gorm"
	"time"
)

type Validation struct {
	gorm.Model
	VAT         string
	CountryCode string
	TradeName   string
	IsVATValid  bool
	CheckedAt   time.Time
}
