package handler

import (
	"net/http"
	"github.com/Marcello1446/vatify/internal/repository"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)


type Request struct {
	VAT         string `json:"vat" validate:"required,numeric,min=8,max=13"`
	CountryCode string `json:"country_code" validate:"required,len=2,alpha,uppercase,oneof=AT BE BG CY CZ DE DK EE EL ES FI FR HR HU IE IT LT LU LV MT NL PL PT RO SE SL SK XI"`
}

type UsageResponse struct {
	Plan         string `json:"plan"`
	RequestsUsed int    `json:"requests_used"`
	RequestLimit int    `json:"request_limit"`
}

type Response struct {
	VAT         string `json:"vat"`
	CountryCode string `json:"country_code"`
	TradeName   string `json:"trade_name"`
	Address     string `json:"address"`
	IsValid     bool   `json:"is_valid"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}


// ValidateHandler Gets VAT and Country Code
// @Summary      Gets VAT and Country Code
// @Description  Gets VAT and Country Code
// @Accept       json
// @Produce      json
// @Param        body   body     Request  true "VAT and country code"
// @Success      200  {object}  Response
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /validate [post]
func ValidateHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		
		k := c.Get("api_key").(string)
		var req Request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "failed to bind data",
			})
		}

		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to validate your request. Country Code either VAT number has incorrect format",
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

// SeeUsage Gets already used requests, plan and limit
// @Summary      Gets already used requests, plan and limit
// @Description  Gets already used requests, plan and limit
// @Accept       json
// @Produce      json
// @Success      200  {object}  UsageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /usage [get]
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