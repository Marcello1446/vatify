package main

import (
	"github.com/Marcello1446/vatify/config"
	_ "github.com/Marcello1446/vatify/docs"
	"github.com/Marcello1446/vatify/internal/database"
	"github.com/Marcello1446/vatify/internal/handler"
	"github.com/Marcello1446/vatify/internal/middleware"
	"github.com/Marcello1446/vatify/internal/model"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
	"log"
	"net/http"
	"time"
)

// @title  VATify API
// @version 1.0
// @description VATify is a reliable EU VAT number validation API. It wraps the official VIES service with caching, audit logging, and tiered API key access.
// @host localhost:8080
// @BasePath /v1

var echo_config = echoMiddleware.RateLimiterConfig{
	Skipper: echoMiddleware.DefaultSkipper,
	Store: echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{Rate: 0.00083, Burst: 3, ExpiresIn: 3 * time.Minute},
	),
	IdentifierExtractor: func(c *echo.Context) (string, error) {
		id := c.RealIP()
		return id, nil
	},
	ErrorHandler: func(c *echo.Context, err error) error {
		return c.JSON(http.StatusForbidden, map[string]string{
			"Error": "Refuses to let you access a webpage or file",
		})
	},
	DenyHandler: func(c *echo.Context, identifier string, err error) error {
		return c.JSON(http.StatusTooManyRequests, map[string]string{
			"Error": "Too many requests. Try again later",
		})
	},
}

var scheduler = cron.New()

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

	scheduler.AddFunc("@daily", func() {
		if err := service.ResetRequests(db); err != nil {
			log.Fatal("Reset error", err)
		}
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/v1/validate", handler.ValidateHandler(db), middleware.AuthAPIKey(db))
	e.POST("/v1/keys", handler.APIKeyHandler(db), echoMiddleware.RateLimiterWithConfig(echo_config))
	e.GET("/v1/usage", handler.SeeUsage(db), middleware.AuthAPIKey(db))
	e.POST("/v1/checkout/sessions", handler.CheckoutHandler(db, *env))
	e.POST("/v1/webhook", handler.WebhookHandler(db, *env))
	e.GET("/*", func(c *echo.Context) error {
		return c.File("index.html")
	})
	e.GET("/success", func(c *echo.Context) error {
		return c.File("success.html")
	})
	e.GET("/cancel", func(c *echo.Context) error { return c.File("cancel.html") })

	e.Static("/img/*", "img")

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "X-API-Key"},
	}))

	scheduler.Start()

	if err := e.Start(":" + env.ServerPort); err != nil {
		log.Println("Failed to start a server:", err)
	}

}
