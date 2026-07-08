package handler

import (
	"context"
	"github.com/Marcello1446/vatify/config"
	"github.com/labstack/echo/v5"
	"github.com/stripe/stripe-go/v86"
	"gorm.io/gorm"
	"net/http"
)

type StripeRequest struct {
	APIKey        string `json:"api_key"`
	Plan          string `json:"plan"`
	SuccessURL    string `json:"success_url"`
	CancelURL     string `json:"cancel_url"`
	CustomerEmail string `json:"email"`
}

func CheckoutHandler(db *gorm.DB, env config.Env) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var stripe_request StripeRequest

		if err := c.Bind(&stripe_request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Failed to bind stripe request data",
			})
		}

		plan := stripe_request.Plan

		price_id := ""
		request_limit := ""

		if plan == "starter" {
			price_id = env.StarterPriceId
		} else if plan == "pro" {
			price_id = env.ProPriceId
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"Error": "Invalid plan",
			})
		}

		if plan == "starter" {
			request_limit = "10,000"
		} else if plan == "pro" {
			request_limit = "100,000"
		}

		sc := stripe.NewClient(env.StripeSecretKey)
		params := &stripe.CheckoutSessionCreateParams{
			CustomerEmail: stripe.String(stripe_request.CustomerEmail),
			SuccessURL:    stripe.String(stripe_request.SuccessURL),
			CancelURL:     stripe.String(stripe_request.CancelURL),
			LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
				{
					Price:    stripe.String(string(price_id)),
					Quantity: stripe.Int64(1),
				},
			},
			Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			Metadata: map[string]string{
				"plan":          plan,
				"request_limit": request_limit,
				"api_key":       stripe_request.APIKey,
			},
		}

		result, err := sc.V1CheckoutSessions.Create(context.TODO(), params)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"Error": "Failed to create checkout session",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"checkout_url": result.URL,
		})
	}
}
