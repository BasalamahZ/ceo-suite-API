package booking

import "errors"

// Followings are the known errors returned from booking.
var (
	// ErrDataNotFound is returned when the desired data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidBookingID is returned when the given booking ID is
	// invalid.
	ErrInvalidBookingID = errors.New("invalid booking id")

	// ErrInvalidStatus is returned when the given status is
	// invalid.
	ErrInvalidStatus = errors.New("invalid status")
)
