package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/Marcello1446/vatify/internal/model"
	"github.com/Marcello1446/vatify/internal/repository"
	"gorm.io/gorm"
	"time"
)

func GenerateAPIKey(db *gorm.DB, email string, plan string) (string, error) {

	_, err := repository.FindAPIKeyByEmail(db, email)

	if err == nil {
		return "", errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	var request_limit int

	if plan == "free" {
		request_limit = 20
	} else if plan == "starter" {
		request_limit = 10000
	} else if plan == "pro" {
		request_limit = 100000
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	plainKey := hex.EncodeToString(key)
	hashedKey := HashAPIKey(plainKey)

	api_key_instance := &model.APIKey{
		StripeCustomerID:  "",
		Key:               hashedKey,
		IsActive:          true,
		Owner:             email,
		Plan:              plan,
		RequestLimit:      request_limit,
		RequestsUsed:      0,
		UsageReset:        time.Now().AddDate(0, 1, 0),
		CancelAtPeriodEnd: false,
		CancelAt:          nil,
	}

	if err := repository.SaveAPIKey(db, api_key_instance); err != nil {
		return "", err
	}

	return plainKey, nil
}

func ResetRequests(db *gorm.DB) error {
	api_keys, err := repository.FindAllAPIKeysUsageReset(db)

	if err != nil {
		return err
	}

	for _, v := range api_keys {
		err := repository.ResetRequestsUsed(db, &v)

		if err != nil {
			return err
		}
	}

	return nil
}

func UpgradePlan(db *gorm.DB, email, plan, customer_id string) error {
	api_key_instance, err := repository.FindAPIKeyByEmail(db, email)

	if err != nil {
		return err
	}

	var request_limit int

	if plan == "starter" {
		request_limit = 10000
	} else if plan == "pro" {
		request_limit = 100000
	} else {
		return errors.New("Plan is not valid - choose pro or starter")
	}

	api_key_instance.StripeCustomerID = customer_id
	api_key_instance.Plan = plan
	api_key_instance.RequestLimit = request_limit
	api_key_instance.RequestsUsed = 0
	api_key_instance.UsageReset = time.Now().AddDate(0, 1, 0)

	if res := db.Save(&api_key_instance); res.Error != nil {
		return res.Error
	}

	return nil
}

func DowngradePlan(db *gorm.DB, customer_id string) error {
	api_key_instance, err := repository.FindAPIKeyByStripeCustomerID(db, customer_id)

	if err != nil {
		return err
	}

	request_limit := 20

	api_key_instance.Plan = "free"
	api_key_instance.RequestLimit = request_limit
	api_key_instance.UsageReset = time.Now().AddDate(0, 1, 0)
	api_key_instance.CancelAtPeriodEnd = false
	api_key_instance.CancelAt = nil

	if res := db.Save(&api_key_instance); res.Error != nil {
		return res.Error
	}

	return nil
}

func UpdateCancellationStatus(db *gorm.DB, customerID string, cancelAtPeriodEnd bool, cancelAt *time.Time) error {
	api_key_instance, err := repository.FindAPIKeyByStripeCustomerID(db, customerID)

	if err != nil {
		return err
	}

	if cancelAtPeriodEnd {
		api_key_instance.CancelAtPeriodEnd = cancelAtPeriodEnd
		api_key_instance.CancelAt = cancelAt
	} else {
		api_key_instance.CancelAtPeriodEnd = false
		api_key_instance.CancelAt = nil
	}

	if res := db.Save(&api_key_instance); res.Error != nil {
		return res.Error
	}

	return nil
}
