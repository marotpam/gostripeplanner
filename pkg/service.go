package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
)

type Service struct {
	stripeEnvironments []*Environment
}

func NewService(envs []*Environment) *Service {
	return &Service{envs}
}

func (s *Service) CopyAllProducts(src, dest string) error {
	srcEnv, found := s.findEnvironment(src)
	if !found {
		return errors.New(fmt.Sprintf("Src environment %s was not configured", src))
	}

	destEnv, found := s.findEnvironment(dest)
	if !found {
		return errors.New(fmt.Sprintf("Dest environment %s was not configured", src))
	}

	srcProductsSvc := NewProductsService(srcEnv)
	productsInSrc, err := srcProductsSvc.All()
	if err != nil {
		return err
	}

	destProductsSvc := NewProductsService(destEnv)
	srcPlansSvc := NewPlansService(srcEnv)
	destPlansSvc := NewPlansService(destEnv)
	for _, p := range productsInSrc {
		prodParams := &stripe.ProductParams{
			Attrs:               p.Attrs,
			Caption:             p.Caption,
			DeactivateOn:        p.DeactivateOn,
			Desc:                p.Desc,
			ID:                  p.ID,
			Images:              p.Images,
			Name:                p.Name,
			StatementDescriptor: p.StatementDescriptor,
			Type:                p.Type,
			URL:                 p.URL,
		}
		if p.Type == stripe.ProductTypeGood {
			prodParams.Active = &p.Active
			prodParams.PackageDimensions = p.PackageDimensions
			prodParams.Shippable = &p.Shippable
		}

		destProduct, err := destProductsSvc.Create(prodParams)
		if err != nil {
			return err
		}

		srcPlans, err := srcPlansSvc.FindForProduct(destProduct.ID)
		if err != nil {
			return err
		}

		for _, sp := range srcPlans {
			var planTiers []*stripe.PlanTierParams
			for _, t := range sp.Tiers {
				planTiers = append(planTiers, &stripe.PlanTierParams{Amount: t.Amount, UpTo: t.UpTo})
			}

			var transformUsage *stripe.PlanTransformUsageParams
			if sp.TransformUsage != nil {
				transformUsage = &stripe.PlanTransformUsageParams{
					DivideBy: sp.TransformUsage.DivideBy,
					Round:    sp.TransformUsage.Round,
				}
			}

			planParams := &stripe.PlanParams{
				Amount:         sp.Amount,
				AmountZero:     sp.Amount == 0,
				BillingScheme:  sp.BillingScheme,
				Currency:       sp.Currency,
				ID:             sp.ID,
				Interval:       sp.Interval,
				IntervalCount:  sp.IntervalCount,
				Nickname:       sp.Nickname,
				ProductID:      &p.ID,
				Tiers:          planTiers,
				TiersMode:      sp.TiersMode,
				TransformUsage: transformUsage,
				TrialPeriod:    sp.TrialPeriod,
				UsageType:      sp.UsageType,
			}

			_, err := destPlansSvc.Create(planParams)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) findEnvironment(name string) (*Environment, bool) {
	for _, e := range s.stripeEnvironments {
		if e.Name == name {
			return e, true
		}
	}

	return nil, false
}
