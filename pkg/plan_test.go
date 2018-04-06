package pkg_test

import (
	"github.com/marotpam/gostripeplanner/pkg"
	"github.com/stripe/stripe-go"
	"testing"
)

func TestItCanCreatePlans(t *testing.T) {
	svc := pkg.NewPlansService(pkg.NewMockedEnvironment("env", "http://0.0.0.0:8420/v1"))

	pID := "prod_123"
	p, err := svc.Create(&stripe.PlanParams{
		Currency:  "eur",
		Interval:  "week",
		ProductID: &pID,
		Amount:    5000,
	})

	if err != nil {
		t.Errorf("Error creating plan: %s", err)
	}

	if p == nil {
		t.Error("Product should not be nil after creation")
	}
}
