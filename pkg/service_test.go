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
	srcProduct, err := srcProductsSvc.Create(&stripe.ProductParams{Name: "product from source", Type: "service"})
	if err != nil {
		t.Fatalf("Error creating product in source: %s", err)
	}

	srcPlansSvc := pkg.NewPlansService(srcEnvironment)
	srcPlan, err := srcPlansSvc.Create(&stripe.PlanParams{
		Amount:    5000,
		Currency:  "eur",
		ProductID: &srcProduct.ID,
		Interval:  "month",
	})
	if err != nil {
		t.Fatalf("Error creating plan in src: %s", err)
	}

	svc := pkg.NewService([]*pkg.Environment{srcEnvironment, destEnvironment})

	err = svc.CopyAllProducts(srcEnvironment.Name, destEnvironment.Name)
	if err != nil {
		t.Fatalf("Error copying all products: %s", err)
	}

	destProductsSvc := pkg.NewProductsService(destEnvironment)
	destProduct, err := destProductsSvc.GetByID(srcProduct.ID)
	if err != nil {
		t.Errorf("Error retrieving product copied from src: %s", err)
	}

	if !reflect.DeepEqual(srcProduct, destProduct) {
		t.Errorf("Product should be the same in source environment than dest environment")
	}

	destPlansSvc := pkg.NewPlansService(destEnvironment)
	destPlan, err := destPlansSvc.GetById(srcPlan.ID)
	if err != nil {
		t.Errorf("Error retrieving plan copied from src: %s", err)
	}

	if !reflect.DeepEqual(srcPlan, destPlan) {
		t.Errorf("Plan should be the same in source environment than dest environment")
	}
}
