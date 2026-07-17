package handler

import (
	"fmt"
	"github.com/Marcello1446/vatify/internal/repository"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"sync"
)

type BulkRequestItem struct {
	SerialNumber int    `json:"-"`
	VAT          string `json:"vat" validate:"required,min=8,max=13"`
	CountryCode  string `json:"country_code" validate:"required,len=2,alpha,uppercase,oneof=AT BE BG CY CZ DE DK EE EL ES FI FR HR HU IE IT LT LU LV MT NL PL PT RO SE SL SK XI"`
}

type BulkRequest struct {
	RequestItems []BulkRequestItem `json:"request_items" validate:"required,min=2,max=100,dive"`
}
type BulkResultItem struct {
	SerialNumber int    `json:"-"`
	VAT          string `json:"vat"`
	CountryCode  string `json:"country_code"`
	IsValid      bool   `json:"is_valid,omitempty"`
	TradeName    string `json:"trade_name,omitempty"`
	Address      string `json:"address,omitempty"`
	Error        string `json:"error,omitempty"`
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
			work := make(chan BulkRequestItem, len(request.RequestItems))
			result := make(chan BulkResultItem, len(request.RequestItems))
			var wg sync.WaitGroup

			go func() {
				for i, item := range request.RequestItems {
					item.SerialNumber = i + 1
					work <- item
				}
				close(work)
			}()

			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for item := range work {
						func() {
							defer func() {
								if r := recover(); r != nil {
									result <- BulkResultItem{VAT: item.VAT, CountryCode: item.CountryCode, Error: fmt.Sprintf("internal error: %v", r)}
								}
							}()

							val, err := service.ValidateVATWithoutIncrementing(db, item.VAT, item.CountryCode)

							if err != nil {
								result <- BulkResultItem{
									SerialNumber: item.SerialNumber,
									VAT:          item.VAT,
									CountryCode:  item.CountryCode,
									IsValid:      false,
									TradeName:    "",
									Address:      "",
									Error:        err.Error(),
								}
							} else {
								result <- BulkResultItem{
									SerialNumber: item.SerialNumber,
									VAT:          val.VAT,
									CountryCode:  val.CountryCode,
									IsValid:      val.IsVATValid,
									TradeName:    val.TradeName,
									Address:      val.Address,
									Error:        "",
								}
							}

							go func() {
								wg.Wait()
								close(result)
							}()

							for r := range result {
								response.ResponseItems = append(response.ResponseItems, r)
								fmt.Println(len(response.ResponseItems))
							}

							sort.Slice(response.ResponseItems, func(i, j int) bool {
								return response.ResponseItems[i].SerialNumber < response.ResponseItems[j].SerialNumber
							})

							fmt.Println("Res inside the function", response.ResponseItems)

							if err := repository.IncreaseMultipleUsedRequests(db, k, len(request.RequestItems)); err != nil {
								fmt.Println("Failed to increase usage count for key:", k, "error:", err)
							}

						}()
						fmt.Println("Res inside the function", response.ResponseItems)

					}
					fmt.Println("Res inside the function", response.ResponseItems)

				}()
				fmt.Println("Res inside the function", response.ResponseItems)
			}
			fmt.Println("Res inside the function", response.ResponseItems)
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Insufficient credentials to make such a big request",
			})
		}

		fmt.Println("Res outside: ", response.ResponseItems)
		return c.JSON(http.StatusOK, &response.ResponseItems)
	}
}
