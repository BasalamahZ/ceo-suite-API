package service

import (
	"context"

	"github.com/ceo-suite/internal/booking"
)

// GetBookingByID returns a booking with the given booking ID.
func (s *service) GetBookingByID(ctx context.Context, id int64) (booking.Booking, error) {
	// validate id
	if id <= 0 {
		return booking.Booking{}, booking.ErrInvalidBookingID
	}

	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return booking.Booking{}, err
	}

	// get a product from postgre
	result, err := pgStoreClient.GetBookingByID(ctx, id)
	if err != nil {
		return booking.Booking{}, err
	}

	return result, nil
}

// GetAllBooking returns list of booking that satisfy the given
// filter.
func (s *service) GetAllBooking(ctx context.Context, filter booking.GetBookingFilter) ([]booking.Booking, error) {
	// get pg store client
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get all booking from postgre
	result, err := pgStoreClient.GetAllBooking(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
