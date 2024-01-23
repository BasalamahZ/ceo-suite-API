package service

import (
	"context"

	"github.com/ceo-suite/internal/product"
)

// GetProductByID returns a product with the given product ID.
func (s *service) GetProductByID(ctx context.Context, id int64) (product.Product, error) {
	// validate id
	if id <= 0 {
		return product.Product{}, product.ErrInvalidProductID
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return product.Product{}, err
	}

	// get a product from postgre
	result, err := pgStoreClient.GetProductByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}

	return result, nil
}

// GetProducts returns list of products that satisfy the given
// filter.
func (s *service) GetProducts(ctx context.Context, filter product.GetProductsFilter) ([]product.Product, error) {
	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all product from postgre
	result, err := pgStoreClient.GetProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
