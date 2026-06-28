package handler

import (
	"net/http"

	"github.com/Marcello1446/vatify/internal/repository"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Request struct {
	VAT         string `json:"vat"`
	CountryCode string `json:"country_code"`
}

type Response struct {
	VAT         string `json:"vat"`
	CountryCode string `json:"country_code"`
	TradeName   string `json:"trade_name"`
	Address     string `json:"address"`
	IsValid     bool   `json:"is_valid"`
}

func ValidateHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		k := c.Get("api_key").(string)
		var req *Request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "failed to bind data",
			})
		}

		val, err := service.ValidateVAT(k, db, req.VAT, req.CountryCode)

		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"Error": "failed to call validation",
			})
		}

		return c.JSON(http.StatusOK, &Response{
			VAT:         val.VAT,
			CountryCode: val.CountryCode,
			TradeName:   val.TradeName,
			Address:     val.Address,
			IsValid:     val.IsVATValid,
		})
	}
}


func SeeUsage(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		k := c.Get("api_key").(string)

		api_key_instance, err := repository.FindAPIKey(db, k)
		
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"Error": "Could not find api key",
			})
		}


		return c.JSON(http.StatusOK, map[string]any{
			"requests_used": api_key_instance.RequestsUsed,
			"request_limit": api_key_instance.RequestLimit,
			"plan": api_key_instance.Plan,
		})
	}
}