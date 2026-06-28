package main

import (
	"github.com/Marcello1446/vatify/config"
	"github.com/Marcello1446/vatify/internal/database"
	"github.com/Marcello1446/vatify/internal/handler"
	"github.com/Marcello1446/vatify/internal/model"
	"github.com/Marcello1446/vatify/internal/middleware"
	"github.com/labstack/echo/v5"
	"github.com/go-playground/validator/v10"
	"log"
)


type CustomValidator struct {
  validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
  if err := cv.validator.Struct(i); err != nil {
    return echo.ErrBadRequest.Wrap(err)
  }
  return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

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
	e.GET("/v1/usage", handler.SeeUsage(db), middleware.AuthAPIKey(db))

	if err := e.Start(":" + env.ServerPort); err != nil {
		log.Fatal("Failed to start a server:", err)
	}
}
