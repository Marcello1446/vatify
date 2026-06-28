package repository

import (
	"github.com/Marcello1446/vatify/internal/model"
	"gorm.io/gorm"
	"time"
)

func SaveAPIKey(db *gorm.DB, a *model.APIKey) error {
	res := db.Create(&a)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func FindAPIKey(db *gorm.DB, key string) (*model.APIKey, error) {
	var api_key *model.APIKey

	if res := db.Where("key = ?", key).First(&api_key); res.Error != nil {
		return nil, res.Error
	}

	return api_key, nil
}


func IncreaseUsedRequests(db *gorm.DB, key string) error {
	var api_key *model.APIKey

	if res := db.Where("key = ?", key).First(&api_key); res.Error != nil {
		return res.Error
	}

	api_key.RequestsUsed += 1

	db.Save(api_key)

	return nil
}


func FindAllAPIKeysUsageReset(db *gorm.DB) ([]model.APIKey ,error) {
	var api_keys []model.APIKey

	if res := db.Where("usage_reset < ?", time.Now()).Find(&api_keys); res.Error != nil {
		return nil, res.Error
	}

	return api_keys, nil
}


func ResetRequestsUsed(db *gorm.DB, api_key *model.APIKey) error {
	api_key.RequestsUsed = 0
	api_key.UsageReset = time.Now().AddDate(0, 1, 0)


	if res := db.Save(&api_key); res.Error != nil {
		return res.Error
	}

	return nil
} 