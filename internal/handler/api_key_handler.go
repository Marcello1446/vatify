package handler

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"net/http"
)

type Login struct {
	Email string `json:"email"`
	Plan  string `json:"plan"`
}

func APIKeyHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var login *Login

		if err := c.Bind(&login); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to bind data",
			})
		}

		return c.JSON(http.StatusOK, login)
	}
}
