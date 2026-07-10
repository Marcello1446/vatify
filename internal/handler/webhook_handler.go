package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Marcello1446/vatify/config"
	"github.com/Marcello1446/vatify/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func WebhookHandler(db *gorm.DB, env config.Env) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var checkout_session *stripe.CheckoutSession
		var subscription_object *stripe.Subscription
		var subscription_customer_id string

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

		if event.Type == "checkout.session.completed" {
			raw_data := event.Data.Raw
			errr := json.Unmarshal(raw_data, &checkout_session)
			customer_id := checkout_session.Customer.ID
			subscription_customer_id = customer_id

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

			if err := service.UpgradePlan(db, checkout_session.CustomerDetails.Email, checkout_session.Metadata["plan"], customer_id); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"Error": "Failed to upgrade plan",
				})
			}

			if checkout_session.Metadata["api_key"] == "existing" {
				if err := service.SendEmailWithUpgradedAPIKey(env, checkout_session.CustomerDetails.Email, checkout_session.Metadata["request_limit"], checkout_session.Metadata["plan"]); err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{
						"Error": "Failed to send email with updated api key",
					})
				}
			} else {
				if err := service.SendEmailWithNewAPIKey(env, checkout_session.Metadata["request_limit"], checkout_session.CustomerDetails.Email, checkout_session.Metadata["api_key"], checkout_session.Metadata["plan"]); err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{
						"Error": "Failed to send email with new api key",
					})
				}
			}

			return c.JSON(http.StatusOK, map[string]string{
				"Message": "Your plan was successfully updated!",
			})
		} else if event.Type == "customer.subscription.deleted" {
			raw_data := event.Data.Raw
			errr := json.Unmarshal(raw_data, &subscription_object)

			if errr != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"Error": "Failed to unmarshal data",
				})
			}

			fmt.Println("Stripe customer id: ", subscription_customer_id)
			if err := service.DowngradePlan(db, subscription_customer_id); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"Error": "Failed to cancel subscription immediately",
				})
			}

			return c.JSON(http.StatusOK, map[string]string{
				"Message": "Your plan was downgraded successfully/ You have cancelled your subscription successfully",
			})

		} else if event.Type == "customer.subscription.updated" {

		}

		return errors.New("Event type does not match our listed ones")
	}
}
