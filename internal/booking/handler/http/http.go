package http

import (
	"errors"
	"net/http"

	"github.com/ceo-suite/internal/booking"
	"github.com/ceo-suite/internal/user"
	"github.com/gorilla/mux"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// dateFormat denotes the standard date format used in
// booking HTTP request and response.
var dateFormat = "02/01/2006"

// timeFormat denotes the standard time format used in
// booking HTTP request and response.
var timeFormat = "3:04 PM"


// Handler contains booking HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	booking  booking.Service
	client   user.Service
}

// handler is the HTTP handler wrapper.
type handler struct {
	h        http.Handler
	identity HandlerIdentity
}

// HandlerIdentity denotes the identity of an HTTP hanlder.
type HandlerIdentity struct {
	Name string
	URL  string
}

// Followings are the known HTTP handler identities
var (
	// HandlerBookings denotes HTTP handler to interact
	// with booking
	HandlerBookings = HandlerIdentity{
		Name: "bookings",
		URL:  "/v1/booking",
	}

	// HandlerBooking denotes HTTP handler to interact
	// with a booking.
	HandlerBooking = HandlerIdentity{
		Name: "booking",
		URL:  "/v1/booking/{id}",
	}
)

// New creates a new Handler.
func New(booking booking.Service, client user.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		booking:  booking,
		client:   client,
	}

	// apply identity
	for _, identity := range identities {
		if h.handlers == nil {
			h.handlers = map[string]*handler{}
		}

		h.handlers[identity.Name] = &handler{
			identity: identity,
		}

		handler, err := h.createHTTPHandler(identity.Name)
		if err != nil {
			return nil, err
		}

		h.handlers[identity.Name].h = handler
	}

	return h, nil
}

// createHTTPHandler creates a new HTTP handler that
// implements http.Handler.
func (h *Handler) createHTTPHandler(configName string) (http.Handler, error) {
	var httpHandler http.Handler
	switch configName {
	case HandlerBookings.Name:
		httpHandler = &bookingsHandler{
			booking: h.booking,
			client:  h.client,
		}
	case HandlerBooking.Name:
		httpHandler = &bookingHandler{
			booking: h.booking,
			client:  h.client,
		}
	default:
		return httpHandler, errUnknownConfig
	}
	return httpHandler, nil
}

// Start starts all HTTP handlers.
func (h *Handler) Start(multiplexer *mux.Router) error {
	for _, handler := range h.handlers {
		multiplexer.Handle(handler.identity.URL, handler.h)
	}
	return nil
}

type bookingHTTP struct {
	ID        *int64  `json:"id"`
	UserID    *int64  `json:"user_id"`
	ProductID *int64  `json:"product_id"`
	Date      *string `json:"date"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Status    *string `json:"status"`
	Price     *int64  `json:"price"`
}
