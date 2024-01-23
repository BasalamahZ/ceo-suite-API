package product

import "errors"

// Followings are the known errors returned from product.
var (
	// ErrDataNotFound is returned when the desired data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidProductID is returned when the given product ID is
	// invalid.
	ErrInvalidProductID = errors.New("invalid product id")
)
