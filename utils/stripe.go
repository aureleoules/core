package utils

import (
	"github.com/stripe/stripe-go/client"
)

var StripeClient client.API

func InitStripe() {
	config := GetConfig()
	StripeClient.Init(config.StripeKey, nil)
}
