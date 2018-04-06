package pkg_test

import (
	"github.com/marotpam/gostripeplanner/pkg"
	"github.com/stripe/stripe-go"
	"reflect"
	"testing"
)

func TestItCanCopyAllProductsBetweenEnvironments(t *testing.T) {
	srcEnvironment := pkg.NewMockedEnvironment("src", "http://0.0.0.0:8420/v1")
	destEnvironment := pkg.NewMockedEnvironment("dest", "http://0.0.0.0:8421/v1")

	srcProductsSvc := pkg.NewProductsService(srcEnvironment)
	sp, err := srcProductsSvc.Create(&stripe.ProductParams{Name: "product from source", Type: "service"})
	if err != nil {
		t.Fatalf("Error creating product in source: %s", err)
	}

	svc := pkg.NewService([]*pkg.Environment{srcEnvironment, destEnvironment})

	err = svc.CopyAllProducts(srcEnvironment.Name, destEnvironment.Name)
	if err != nil {
		t.Fatalf("Error copying all products: %s", err)
	}

	destProductsSvc := pkg.NewProductsService(destEnvironment)
	dp, err := destProductsSvc.GetByID(sp.ID)
	if err != nil {
		t.Errorf("Error retrieving product copied from src: %s", err)
	}

	if !reflect.DeepEqual(sp, dp) {
		t.Errorf("Product should be the same in source environment than dest environment")
	}
}
