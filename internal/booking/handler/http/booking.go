package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ceo-suite/global/helper"
	"github.com/ceo-suite/internal/booking"
	"github.com/ceo-suite/internal/user"
	"github.com/gorilla/mux"
)

type bookingHandler struct {
	booking booking.Service
	client  user.Service
}

func (h *bookingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Booking HTTP][bookingHandler] Failed to parse booking ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidBookingID.Error()})
		return
	}

	// handle based on HTTP request method
	switch r.Method {
	case http.MethodGet:
		h.handleGetBookingByID(w, r, bookingID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *bookingHandler) handleGetBookingByID(w http.ResponseWriter, r *http.Request, bookingID int64) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 1000*time.Millisecond)
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
			log.Printf("[Booking HTTP][handleGetBookingByID] Failed to get product by ID. bookingID: %d, Err: %s\n", bookingID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan booking.Booking, 1)
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
		err = checkAccessToken(ctx, h.client, token, "handleGetBookingByID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.booking.GetBookingByID(ctx, bookingID)
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
				log.Printf("[Booking HTTP][handleGetBookingByID] Internal error from GetBookingByID. bookingID: %d. Err: %s\n", bookingID, err.Error())
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
		// format booking
		var b bookingHTTP
		b, err = formatBooking(res)
		if err != nil {
			return
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: b,
		})
	}
}
