package main

import (
	"github.com/Marcello1446/vatify/config"
	"github.com/Marcello1446/vatify/internal/database"
	"github.com/Marcello1446/vatify/internal/handler"
	"github.com/Marcello1446/vatify/internal/model"
	"github.com/Marcello1446/vatify/internal/middleware"
	"github.com/labstack/echo/v5"
	"log"
)

func main() {
	e := echo.New()
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	db, err := database.InitDB(env)
	if err != nil {
		log.Fatal("Failed to launch database:", err)
	}

	er := db.AutoMigrate(&model.APIKey{}, &model.Validation{})
	if er != nil {
		log.Fatal("Failed to migrate:", er)
	}

	e.POST("/v1/validate", handler.ValidateHandler(db), middleware.AuthAPIKey(db))
	e.POST("/v1/keys", handler.APIKeyHandler(db))

	if err := e.Start(":" + env.ServerPort); err != nil {
		log.Fatal("Failed to start a server:", err)
	}
}
