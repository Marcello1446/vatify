package service

import (
	"errors"
	"github.com/Marcello1446/vatify/internal/model"
	"github.com/Marcello1446/vatify/internal/repository"
	"github.com/Marcello1446/vatify/internal/vies"
	"gorm.io/gorm"
	"time"
)

func ValidateVAT(db *gorm.DB, vat string, countryCode string) (*model.Validation, error) {
	validation, err := repository.FindValidation(db, vat, countryCode)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		res, err := vies.Check(countryCode, vat)

		if err != nil {
			return nil, err
		}

		v := &model.Validation{
			VAT:         res.CleanedVAT,
			CountryCode: res.CountryCode,
			TradeName:   res.TradeName,
			IsVATValid:  res.IsValid,
			CheckedAt:   time.Now(),
		}

		erra := repository.SaveValidation(db, v)

		if erra != nil {
			return nil, erra
		}

		return v, nil

	} else if err != nil {
		return nil, err
	}

	elapsed := time.Since(validation.CheckedAt)

	if elapsed > 24*time.Hour {
		res, err := vies.Check(validation.CountryCode, validation.VAT)

		if err != nil {
			return nil, err
		}

		v := &model.Validation{
			VAT:         res.CleanedVAT,
			CountryCode: res.CountryCode,
			TradeName:   res.TradeName,
			IsVATValid:  res.IsValid,
			CheckedAt:   time.Now(),
		}

		erra := repository.SaveValidation(db, v)

		if erra != nil {
			return nil, erra
		}

		return v, nil

	} else {
		return validation, nil
	}

	return validation, nil
}
