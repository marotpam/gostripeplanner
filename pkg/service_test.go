package pkg_test

import (
	"testing"
	"github.com/marotpam/gostripeplanner/pkg"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go"
	"net/http"
	"time"
	"reflect"
)

func TestItCanRetrieveAllProducts(t *testing.T) {
	svc := pkg.NewService(localStripeClient())

	productsBefore, err := svc.AllProducts()
	if err != nil {
		t.Error(err)
	}

	_, err = svc.AddProduct(&stripe.ProductParams{Name: "new product", Type: "service"})
	if err != nil {
		t.Fatalf("Cannot create stripe product: %s", err)
	}

	productsAfter, err := svc.AllProducts()
	if err != nil {
		t.Error(err)
	}

	if len(productsAfter) != len(productsBefore) + 1 {
		t.Error("List of products after should have increased by 1")
	}
}

func TestProductCreation(t *testing.T) {
	svc := pkg.NewService(localStripeClient())

	newProduct, err := svc.AddProduct(&stripe.ProductParams{Name: "new product", Type: "service"})
	if err != nil {
		t.Fatalf("Cannot create stripe product: %s", err)
	}

	retrievedProduct, err := svc.GetProductByID(newProduct.ID)
	if err != nil {
		t.Fatalf("Error retrieving created product: %s", err)
	}

	if !reflect.DeepEqual(newProduct, retrievedProduct) {
		t.Error("Retrieved product is different than the one created")
	}
}

func localStripeClient() *client.API {
	stripeClient := client.API{}
	c := http.Client{Timeout: 30 * time.Second}
	stripeClient.Init("sk_123", &stripe.Backends{
		API: &stripe.BackendConfiguration{
			stripe.APIBackend,
			"http://0.0.0.0:8420/v1",
			&c,
		},
	})
	return &stripeClient
}
