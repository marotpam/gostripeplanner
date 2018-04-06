package pkg_test

import (
	"fmt"
	"github.com/marotpam/gostripeplanner/pkg"
	"github.com/stripe/stripe-go"
	"reflect"
	"testing"
)

func TestItCanRetrieveAllProducts(t *testing.T) {
	ps := pkg.NewProductsService(localStripeEnvironment())

	productsBefore, err := ps.All()
	if err != nil {
		t.Error(err)
	}

	_, err = ps.Create(&stripe.ProductParams{Name: "new product", Type: "service"})
	if err != nil {
		t.Fatalf("Cannot create stripe product: %s", err)
	}

	productsAfter, err := ps.All()
	if err != nil {
		t.Error(err)
	}

	if len(productsAfter) != len(productsBefore)+1 {
		t.Error("List of products after should have increased by 1")
	}
}

func TestProductCreation(t *testing.T) {
	ps := pkg.NewProductsService(localStripeEnvironment())

	newProduct, err := ps.Create(&stripe.ProductParams{Name: "new product", Type: "service"})
	if err != nil {
		t.Fatalf("Cannot create stripe product: %s", err)
	}

	retrievedProduct, err := ps.GetByID(newProduct.ID)
	if err != nil {
		t.Fatalf("Error retrieving created product: %s", err)
	}

	if !reflect.DeepEqual(newProduct, retrievedProduct) {
		t.Error("Retrieved product is different than the one created")
	}
}

func TestItCanPurgeTheProductCatalog(t *testing.T) {
	ps := pkg.NewProductsService(localStripeEnvironment())

	for i := 0; i < 10; i++ {
		_, err := ps.Create(&stripe.ProductParams{Name: fmt.Sprintf("product with id %s", i), Type: "service"})
		if err != nil {
			t.Fatalf("Cannot create stripe product: %s", err)
		}
	}

	if err := ps.Purge(); err != nil {
		t.Errorf("Error purging products: %s", err)
	}

	products, err := ps.All()
	if err != nil {
		t.Errorf("Error retrieving products after purge: %s", err)
	}

	if len(products) != 0 {
		t.Errorf("Should have 0 products after purging, %d found", len(products))
	}
}

func localStripeEnvironment() *pkg.Environment {
	return pkg.NewMockedEnvironment("mock", "http://0.0.0.0:8420/v1")
}
