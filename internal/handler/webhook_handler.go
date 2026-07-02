package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"fmt"
	"github.com/Marcello1446/vatify/config"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
	"gorm.io/gorm"
)



func WebhookHandler(db *gorm.DB, env config.Env) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var checkout_session *stripe.CheckoutSession

		raw_body, err := io.ReadAll(c.Request().Body)
		

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to read checkout session data",
			})
		}

		stripe_header := c.Request().Header.Get("Stripe-Signature")

	
		event, err := webhook.ConstructEventWithOptions(raw_body, stripe_header, env.StripeWebhookSecret, webhook.ConstructEventOptions{
    		IgnoreAPIVersionMismatch: true,
		})

    	
		if err != nil {
			fmt.Println("verify error:", err)
        	return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to verificate webhook signature",
			})
    	}

		if event.Type != "checkout.session.completed" {
    		return c.JSON(http.StatusOK, nil)
		}

		raw_data := event.Data.Raw
		errr := json.Unmarshal(raw_data, &checkout_session)

		if errr != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to unmarshal data",
			})
		}

		if len(checkout_session.Metadata["plan"]) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to read a 'plan' field",
			})
		}

		if err := service.UpgradePlan(db, checkout_session.CustomerDetails.Email, checkout_session.Metadata["plan"]); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to upgrade plan",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"Message": "Your plan was successfully updated!",
		})
	}
}