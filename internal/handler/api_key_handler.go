package handler

import (
	"net/http"

	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Login struct {
	Email string `json:"email"`
	Plan  string `json:"plan"`
}

type ErrorResponses struct {
    Error string `json:"error"`
}


// APIKeyHandler Generates API key based on plan
// @Summary      Generates API key based on plan
// @Description  Generates API key based on plan
// @Accept       json
// @Produce      json
// @Param        body   body    Login  true "Login and email"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  ErrorResponses
// @Failure      404  {object}  ErrorResponses
// @Router       /keys [post]
func APIKeyHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var login Login

		if err := c.Bind(&login); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to bind data",
			})
		}

		api_key, err := service.HashAPIKey(db, login.Email, login.Plan)

		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"Error": "Failed to collect data from service. APIKEY",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"api_key": api_key,
		})
	}
}
