package repository

import (
	"github.com/Marcello1446/vatify/internal/model"
	"gorm.io/gorm"
)

func SaveValidation(db *gorm.DB, v *model.Validation) error {
	res := db.Create(&v)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func FindValidation(db *gorm.DB, vat string, countryCode string) (*model.Validation, error) {
	var validation *model.Validation

	if res := db.Where("vat = ? AND country_code = ?", vat, countryCode).First(&validation); res.Error != nil {
		return nil, res.Error
	}

	return validation, nil
}
