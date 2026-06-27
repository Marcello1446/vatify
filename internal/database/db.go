package database

import (
	"fmt"
	"github.com/Marcello1446/vatify/config"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(env *config.Env) (*gorm.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", env.DbHost, env.DbUser, env.DbPassword, env.DbName, env.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
