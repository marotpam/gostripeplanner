package pkg

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type plansService struct {
	stripeClient *client.API
}

func NewPlansService(env *Environment) *plansService {
	return &plansService{
		env.GetClient(),
	}
}

func (ps *plansService) Create(pp *stripe.PlanParams) (*stripe.Plan, error) {
	return ps.stripeClient.Plans.New(pp)
}

func (ps *plansService) FindForProduct(productID string) ([]*stripe.Plan, error) {
	var plans []*stripe.Plan

	params := &stripe.PlanListParams{}
	params.Filters.AddFilter("product", "", productID)

	i := ps.stripeClient.Plans.List(params)
	for i.Next() {
		plans = append(plans, i.Plan())
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}
