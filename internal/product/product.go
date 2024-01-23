package product

import (
	"context"
	"time"
)

type Service interface {
	// GetProductByID returns a product with the given product ID.
	GetProductByID(ctx context.Context, id int64) (Product, error)

	// GetProducts returns list of products that satisfy the given
	// filter.
	GetProducts(ctx context.Context, filter GetProductsFilter) ([]Product, error)
}

type Product struct {
	ID          int64
	Name        string
	Images      []string
	Location    string
	Date        time.Time
	StartTime   time.Time
	EndTime     time.Time
	Status      Status
	Capacity    int
	Price       int64
	MinCharge   int64
	DailyRate   int64
	Promo       bool
	PromoPrice  int64
	Address     string
	Distance    float32
	Description string
	Latitude    string
	Longitude   string
	Rating      float32
	CreateTime  time.Time
	UpdateTime  time.Time
}

type GetProductsFilter struct {
	Location  string
	StartTime time.Time
	EndTime   time.Time
	Date      time.Time
	Capacity  int
	Promo     bool
}

// Status denotes status of a product.
type Status int

// Followings are the known product status.
const (
	StatusUnknown  Status = 0
	StatusActive   Status = 1
	StatusInactive Status = 2
)

var (
	// StatusList is a list of valid product status.
	StatusList = map[Status]struct{}{
		StatusActive:   {},
		StatusInactive: {},
	}

	// StatusName maps product status to it's string representation.
	StatusName = map[Status]string{
		StatusActive:   "active",
		StatusInactive: "inactive",
	}
)

// Value returns int value of a product status.
func (s Status) Value() int {
	return int(s)
}

// String returns string representaion of a product status.
func (s Status) String() string {
	return StatusName[s]
}
