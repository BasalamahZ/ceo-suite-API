package http

import (
	"errors"

	"github.com/ceo-suite/internal/booking"
)

// Followings are the known errors from Booking HTTP handlers.
var (
	// errDataNotFound is returned when the desired data is
	// not found.
	errDataNotFound = errors.New("DATA_NOT_FOUND")

	// errInternalServer is returned when there is an
	// unexpected error encountered when processing a request.
	errInternalServer = errors.New("INTERNAL_SERVER_ERROR")

	// errMethodNotAllowed is returned when accessing not
	// allowed HTTP method.
	errMethodNotAllowed = errors.New("METHOD_NOT_ALLOWED")

	// errRequestTimeout is returned when processing time has
	// reached the timeout limit.
	errRequestTimeout = errors.New("REQUEST_TIMEOUT")

	// errInvalidBookingID is returned when the given booking ID is
	// invalid.
	errInvalidBookingID = errors.New("INVALID_BOOKING_ID")

	// errInvalidStatus is returned when the given status is
	// invalid.
	errInvalidStatus = errors.New("INVALID_STATUS")

	// errInvalidToken is returned when the given token is
	// invalid.
	errInvalidToken = errors.New("INVALID_TOKEN")

	// errUnauthorizedAccess is returned when the request
	// is unaothorized.
	errUnauthorizedAccess = errors.New("UNAUTHORIZED_ACCESS")
)

var (
	// mapHTTPError maps service error into HTTP error that
	// categorize as bad request error.
	//
	// Internal server error-related should not be mapped here,
	// and the handler should just return `errInternal` as the
	// error instead
	mapHTTPError = map[error]error{
		booking.ErrDataNotFound:     errDataNotFound,
		booking.ErrInvalidBookingID: errInvalidBookingID,
		booking.ErrInvalidStatus:    errInvalidStatus,
	}
)
