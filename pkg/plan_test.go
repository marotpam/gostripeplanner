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

func TestItCanRetrieveAllPlansForAProduct(t *testing.T) {
	svc := pkg.NewPlansService(pkg.NewMockedEnvironment("env", "http://0.0.0.0:8420/v1"))

	ps, err := svc.FindForProduct("prod_123")
	if err != nil {
		t.Errorf("Error getting plans for product: %s", err)
	}

	if len(ps) == 0 {
		t.Errorf("No plan was found")
	}
}
