package pkg

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"net/http"
	"time"
)

type Environment struct {
	Name string
	Key  string
	URL  string
}

func NewMockedEnvironment(name, stripeMockURL string) *Environment {
	return &Environment{
		Name: name,
		Key:  "sk_test_123",
		URL:  stripeMockURL,
	}
}

func (e *Environment) GetClient() *client.API {
	stripeClient := client.API{}
	c := http.Client{Timeout: 30 * time.Second}

	stripeClient.Init(e.Key, &stripe.Backends{
		API: &stripe.BackendConfiguration{
			stripe.APIBackend,
			e.URL,
			&c,
		},
	})

	return &stripeClient
}
