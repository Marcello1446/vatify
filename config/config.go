package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	ServerPort string
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
		ServerPort: os.Getenv("SERVER_PORT")}

	return env, nil
}
