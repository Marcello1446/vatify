package handler

import (
	"github.com/Marcello1446/vatify/internal/repository"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"net/http"
)

type BulkRequestItem struct {
	VAT         string `json:"vat" validate:"required,min=8,max=13"`
	CountryCode string `json:"country_code" validate:"required,len=2,alpha,uppercase,oneof=AT BE BG CY CZ DE DK EE EL ES FI FR HR HU IE IT LT LU LV MT NL PL PT RO SE SL SK XI"`
}

type BulkRequest struct {
	RequestItems []BulkRequestItem `json:"request_items" validate:"required,min=2,max=100,dive"`
}
type BulkResultItem struct {
	VAT         string `json:"vat"`
	CountryCode string `json:"country_code"`
	IsValid     bool   `json:"is_valid,omitempty"`
	TradeName   string `json:"trade_name,omitempty"`
	Address     string `json:"address,omitempty"`
	Error       string `json:"error,omitempty"`
}
type BulkResponse struct {
	ResponseItems []BulkResultItem `json:"response_items"`
}

func BulkValidationHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {

		k := c.Get("api_key").(string)

		api_key_instance, err := repository.FindAPIKey(db, k)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to find a record in a database",
			})
		}

		var request BulkRequest
		var response BulkResponse

		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to bind data",
			})
		}

		if err := c.Validate(&request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to validate your request. Country Code either VAT number has incorrect format or length of your array is greater than 100",
			})
		}

		if api_key_instance.RequestLimit-api_key_instance.RequestsUsed >= len(request.RequestItems) {
			for _, item := range request.RequestItems {
				val, err := service.ValidateVATWithoutIncrementing(db, item.VAT, item.CountryCode)

				if err != nil {
					response.ResponseItems = append(response.ResponseItems, BulkResultItem{
						VAT:         item.VAT,
						CountryCode: item.CountryCode,
						IsValid:     false,
						TradeName:   "",
						Address:     "",
						Error:       err.Error(),
					})
				} else {
					response.ResponseItems = append(response.ResponseItems, BulkResultItem{
						VAT:         val.VAT,
						CountryCode: val.CountryCode,
						IsValid:     val.IsVATValid,
						TradeName:   val.TradeName,
						Address:     val.Address,
						Error:       "",
					})
				}
			}

			if err := repository.IncreaseMultipleUsedRequests(db, k, len(request.RequestItems)); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"Error": "Failed to increase used requests amount",
				})
			}

		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Insufficient credentials to make such a big request",
			})
		}

		return c.JSON(http.StatusOK, response.ResponseItems)
	}
}
