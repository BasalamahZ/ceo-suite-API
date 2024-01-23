package service

import (
	"context"

	"github.com/ceo-suite/internal/booking"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// GetBookingByID returns a booking with the given booking ID.
	GetBookingByID(ctx context.Context, id int64) (booking.Booking, error)

	// GetAllBooking returns list of booking that satisfy the given
	// filter.
	GetAllBooking(ctx context.Context, filter booking.GetBookingFilter) ([]booking.Booking, error)
}
