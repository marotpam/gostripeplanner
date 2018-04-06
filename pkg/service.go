package pkg

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type Service struct {
	stripeClient *client.API
}

func NewService(sc *client.API) *Service {
	return &Service{
		sc,
	}
}

func (s *Service) AllProducts() ([]*stripe.Product, error) {
	var pl []*stripe.Product

	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "100")
	l := s.stripeClient.Products.List(params)

	for l.Next() {
		pl = append(pl, l.Product())
	}

	if err:= l.Err(); err != nil {
		return nil, err
	}

	return pl, nil
}

func (s *Service) AddProduct(pp *stripe.ProductParams) (*stripe.Product, error) {
	return s.stripeClient.Products.New(pp)
}

func (s *Service) GetProductByID(productID string) (*stripe.Product, error) {
	return s.stripeClient.Products.Get(productID, nil)
}