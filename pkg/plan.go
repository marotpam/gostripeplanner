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

func NewWithClient(c *client.API) *plansService {
	return &plansService{c}
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

func (ps *plansService) GetById(planID string) (*stripe.Plan, error) {
	return ps.stripeClient.Plans.Get(planID, nil)
}

func (ps *plansService) DeleteById(planID string) error {
	_, err := ps.stripeClient.Plans.Del(planID, nil)
	return err
}
