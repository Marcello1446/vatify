package handler

import (
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"net/http"
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
		var req *Request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "failed to bind data",
			})
		}

		val, err := service.ValidateVAT(db, req.VAT, req.CountryCode)

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
