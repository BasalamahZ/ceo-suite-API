package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ceo-suite/global/helper"
	"github.com/ceo-suite/internal/booking"
	"github.com/ceo-suite/internal/user"
)

type bookingsHandler struct {
	booking booking.Service
	client  user.Service
}

func (h *bookingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetAllBooking(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *bookingsHandler) handleGetAllBooking(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 2000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[Booking HTTP][handleGetAllBooking] Failed to get all booking. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []booking.Booking, 1)
	errChan := make(chan error, 1)

	go func() {
		// get token from header
		token, err := helper.GetBearerTokenFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errInvalidToken
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.client, token, "handleGetAllBooking")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// parsed filter
		filter, err := parseGetAllBookingFilter(r.URL.Query())
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.booking.GetAllBooking(ctx, filter)
		if err != nil {
			// determine error and status code, by default its internal error
			parsedErr := errInternalServer
			statusCode = http.StatusInternalServerError
			if v, ok := mapHTTPError[err]; ok {
				parsedErr = v
				statusCode = http.StatusBadRequest
			}

			// log the actual error if its internal error
			if statusCode == http.StatusInternalServerError {
				log.Printf("[Booking HTTP][handleGetAllBooking] Internal error from GetAllBooking. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- res
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case res := <-resChan:
		// format each bookings
		bookings := make([]bookingHTTP, 0)
		for _, r := range res {
			var b bookingHTTP
			b, err = formatBooking(r)
			if err != nil {
				return
			}
			bookings = append(bookings, b)
		}

		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: bookings,
		})
	}
}

func parseGetAllBookingFilter(request url.Values) (booking.GetBookingFilter, error) {
	result := booking.GetBookingFilter{}

	var status booking.Status
	if statusStr := request.Get("status"); statusStr != "" {
		parseStatus, err := parseStatus(statusStr)
		if err != nil {
			return result, err
		}
		status = parseStatus
	}

	return booking.GetBookingFilter{
		Status: status,
	}, nil
}
