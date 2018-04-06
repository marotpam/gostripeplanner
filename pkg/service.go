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
		return errors.New(fmt.Sprintf("Src environment %s was not configured", src))
	}

	srcProductsSvc := NewProductsService(srcEnv)
	productsInSrc, err := srcProductsSvc.All()
	if err != nil {
		return err
	}

	destProductsSvc := NewProductsService(destEnv)
	for _, p := range productsInSrc {
		s.copyProduct(p, destProductsSvc)
	}

	return nil
}

func (s *Service) copyProduct(p *stripe.Product, dest *productsService) error {
	params := &stripe.ProductParams{
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
	if p.Type == "good" {
		params.Active = &p.Active
		params.PackageDimensions = p.PackageDimensions
		params.Shippable = &p.Shippable
	}

	_, err := dest.Add(params)
	if err != nil {
		return err
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
