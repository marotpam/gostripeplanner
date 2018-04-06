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

func (ps *productsService) Create(pp *stripe.ProductParams) (*stripe.Product, error) {
	return ps.stripeClient.Products.New(pp)
}

func (ps *productsService) GetByID(productID string) (*stripe.Product, error) {
	return ps.stripeClient.Products.Get(productID, nil)
}

func (ps *productsService) Purge() error {
	products, err := ps.All()
	if err != nil {
		return err
	}

	plansSvc := NewWithClient(ps.stripeClient)
	for _, p := range products {
		plans, err := plansSvc.FindForProduct(p.ID)
		if err != nil {
			return err
		}

		for _, pl := range plans {
			if err := plansSvc.DeleteById(pl.ID); err != nil {
				return err
			}
		}

		if _, err := ps.stripeClient.Products.Del(p.ID, nil); err != nil {
			return err
		}
	}

	return nil
}
