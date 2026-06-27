package repository

import (
	"github.com/Marcello1446/vatify/internal/model"
	"gorm.io/gorm"
)

func SaveAPIKey(db *gorm.DB, a *model.APIKey) error {
	res := db.Create(&a)
	if res.Error != nil {
		return res.Error
	}
	return nil
}