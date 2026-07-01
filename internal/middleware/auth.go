package middleware

import (
	"net/http"

	"github.com/Marcello1446/vatify/internal/repository"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func AuthAPIKey(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
		
		key := c.Request().Header.Get("X-API-Key")
		

		if len(key) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string {
				"Error": "Key length is 0",
			})
		}

		api_key_instance, err := repository.FindAPIKey(db, key)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"Error": "Failed to find key in a database",
			})
		}

		if !api_key_instance.IsActive {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"Error": "Key is not active",
			})
		}

		if api_key_instance.RequestsUsed >= api_key_instance.RequestLimit {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"Error": "Your limit is exceeded",
			})
		}

		c.Set("api_key", key)

		return next(c)
	}
	}

}