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
