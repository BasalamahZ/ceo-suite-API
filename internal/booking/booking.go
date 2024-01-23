package booking

import (
	"context"
	"time"
)

type Service interface {
	// GetBookingByID returns a booking with the given booking ID.
	GetBookingByID(ctx context.Context, id int64) (Booking, error)

	// GetAllBooking returns list of booking that satisfy the given
	// filter.
	GetAllBooking(ctx context.Context, filter GetBookingFilter) ([]Booking, error)
}

type Booking struct {
	ID         int64
	UserID     int64
	ProductID  int64
	Date       time.Time // derived
	StartTime  time.Time // derived
	EndTime    time.Time // derived
	Price      int64     // derived
	Status     Status
	CreateTime time.Time
	UpdateTime time.Time
}

type GetBookingFilter struct {
	Status Status
}

// Status denotes status of a booking.
type Status int

// Followings are the known booking status.
const (
	StatusUnknown   Status = 0
	StatusOngoing   Status = 1
	StatusScheduled Status = 2
	StatusCompleted Status = 3
)

var (
	// StatusList is a list of valid booking status.
	StatusList = map[Status]struct{}{
		StatusOngoing:   {},
		StatusScheduled: {},
		StatusCompleted: {},
	}

	// StatusName maps booking status to it's string representation.
	StatusName = map[Status]string{
		StatusOngoing:   "ongoing",
		StatusScheduled: "scheduled",
		StatusCompleted: "completed",
	}
)

// Value returns int value of a booking status.
func (s Status) Value() int {
	return int(s)
}

// String returns string representaion of a booking status.
func (s Status) String() string {
	return StatusName[s]
}
