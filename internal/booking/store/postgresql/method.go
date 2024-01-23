package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ceo-suite/internal/booking"
	"github.com/jmoiron/sqlx"
)

// GetBookingByID returns a booking with the given booking ID.
func (sc *storeClient) GetBookingByID(ctx context.Context, id int64) (booking.Booking, error) {
	query := fmt.Sprintf(queryGetBooking, "WHERE b.id = $1")
	// query single row
	var bdb bookingDB
	err := sc.q.QueryRowx(query, id).StructScan(&bdb)
	if err != nil {
		if err == sql.ErrNoRows {
			return booking.Booking{}, booking.ErrDataNotFound
		}
		return booking.Booking{}, err
	}

	return bdb.format(), nil
}

// GetAllBooking returns list of booking that satisfy the given
// filter.
func (sc *storeClient) GetAllBooking(ctx context.Context, filter booking.GetBookingFilter) ([]booking.Booking, error) {
	// define variables to custom query
	argsKV := make(map[string]interface{})
	addConditions := make([]string, 0)

	if filter.Status.String() != "" {
		addConditions = append(addConditions, "b.status = :status")
		argsKV["status"] = filter.Status
	}

	// construct strings to custom query
	addCondition := strings.Join(addConditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(addConditions) > 0 {
		addCondition = fmt.Sprintf("WHERE %s", addCondition)
	}

	// construct query
	query := fmt.Sprintf(queryGetBooking, addCondition)

	// prepare query
	query, args, err := sqlx.Named(query, argsKV)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sc.q.Rebind(query)

	// query to database
	rows, err := sc.q.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// read products
	products := make([]booking.Booking, 0)
	for rows.Next() {
		var row bookingDB
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		products = append(products, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
