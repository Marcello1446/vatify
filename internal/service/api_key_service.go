package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"
	"gorm.io/gorm"
	"github.com/Marcello1446/vatify/internal/model"
	"github.com/Marcello1446/vatify/internal/repository"
)

func HashAPIKey(db *gorm.DB, email string, plan string) (string, error) {
	var request_limit int

	if plan == "free" {
		request_limit = 500
	} else if plan == "starter" {
		request_limit = 10000
	} else if plan == "pro" {
		request_limit = 100000
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
    	return "", err
	}
	encoded_string := hex.EncodeToString(key)

	api_key_instance := &model.APIKey{
		Key: encoded_string,
		IsActive: true,
		Owner: email,
		Plan: plan,
		RequestLimit: request_limit,
		RequestsUsed: 0,
		UsageReset: time.Now().AddDate(0, 1, 0),
	}

	errr := repository.SaveAPIKey(db, api_key_instance)

	if errr != nil {
		return "", errr
	}

	return encoded_string, nil
}