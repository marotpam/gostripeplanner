package pkg

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type productsService struct {
	stripeClient *client.API
}

func NewProductsService(env *Environment) *productsService {
	return &productsService{
		env.GetClient(),
	}
}

func (ps *productsService) All() ([]*stripe.Product, error) {
	var pl []*stripe.Product

	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "100")
	l := ps.stripeClient.Products.List(params)

	for l.Next() {
		pl = append(pl, l.Product())
	}

	if err := l.Err(); err != nil {
		return nil, err
	}

	return pl, nil
}

func (ps *productsService) Add(pp *stripe.ProductParams) (*stripe.Product, error) {
	return ps.stripeClient.Products.New(pp)
}

func (ps *productsService) GetByID(productID string) (*stripe.Product, error) {
	return ps.stripeClient.Products.Get(productID, nil)
}
