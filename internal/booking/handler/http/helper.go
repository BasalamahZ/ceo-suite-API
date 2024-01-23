package http

import (
	"github.com/ceo-suite/internal/booking"
)

// formatBooking formats the given booking
// into the respective HTTP-format object.
func formatBooking(p booking.Booking) (bookingHTTP, error) {
	status := p.Status.String()
	date := p.Date.Format(dateFormat)
	startTime := p.StartTime.Format(timeFormat)
	endTime := p.EndTime.Format(timeFormat)

	return bookingHTTP{
		ID:        &p.ID,
		UserID:    &p.UserID,
		ProductID: &p.ProductID,
		Date:      &date,
		StartTime: &startTime,
		EndTime:   &endTime,
		Status:    &status,
		Price:     &p.Price,
	}, nil
}

// parseStatus returns booking.Status from the given string.
func parseStatus(req string) (booking.Status, error) {
	switch req {
	case booking.StatusOngoing.String():
		return booking.StatusOngoing, nil
	case booking.StatusScheduled.String():
		return booking.StatusScheduled, nil
	case booking.StatusCompleted.String():
		return booking.StatusCompleted, nil
	}
	return booking.StatusUnknown, errInvalidStatus
}
