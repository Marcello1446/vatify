package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	DbHost              string
	DbPort              string
	DbUser              string
	DbPassword          string
	DbName              string
	ServerPort          string
	StripeSecretKey     string
	ProPriceId          string
	StarterPriceId      string
	StripeWebhookSecret string

}

func LoadEnv() (*Env, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	env := &Env{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
		StripeSecretKey: os.Getenv("STRIPE_SECRET_KEY"),
		ProPriceId: os.Getenv("STRIPE_PRO_PRICE_ID"),
		StarterPriceId: os.Getenv("STRIPE_STARTER_PRICE_ID"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}

	return env, nil
}
