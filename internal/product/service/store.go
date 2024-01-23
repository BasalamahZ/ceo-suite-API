package service

import (
	"context"

	"github.com/ceo-suite/internal/product"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// GetProductByID returns a product with the given product ID.
	GetProductByID(ctx context.Context, id int64) (product.Product, error)

	// GetProducts returns list of products that satisfy the given
	// filter.
	GetProducts(ctx context.Context, filter product.GetProductsFilter) ([]product.Product, error)
}
