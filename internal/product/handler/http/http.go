package http

import (
	"errors"
	"net/http"

	"github.com/ceo-suite/internal/product"
	"github.com/ceo-suite/internal/user"
	"github.com/gorilla/mux"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// dateFormat denotes the standard date format used in
// product HTTP request and response.
var dateFormat = "02/01/2006"

// timeFormat denotes the standard time format used in
// product HTTP request and response.
var timeFormat = "3:04 PM"

// Handler contains product HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	product  product.Service
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
	// HandlerProducts denotes HTTP handler to interact
	// with products
	HandlerProducts = HandlerIdentity{
		Name: "products",
		URL:  "/v1/products",
	}

	// HandlerProduct denotes HTTP handler to interact
	// with a product.
	HandlerProduct = HandlerIdentity{
		Name: "product",
		URL:  "/v1/products/{id}",
	}
)

// New creates a new Handler.
func New(product product.Service, client user.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		product:  product,
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
	case HandlerProducts.Name:
		httpHandler = &productsHandler{
			product: h.product,
			client:  h.client,
		}
	case HandlerProduct.Name:
		httpHandler = &productHandler{
			product: h.product,
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

type productHTTP struct {
	ID          *int64    `json:"id"`
	Name        *string   `json:"name"`
	Images      *[]string `json:"images"`
	Location    *string   `json:"location"`
	Date        *string   `json:"date"`
	StartTime   *string   `json:"start_time"`
	EndTime     *string   `json:"end_time"`
	Status      *string   `json:"status"`
	Capacity    *int      `json:"capacity"`
	Price       *int64    `json:"price"`
	MinCharge   *int64    `json:"min_charge"`
	DailyRate   *int64    `json:"daily_rate"`
	Promo       *bool     `json:"promo"`
	PromoPrice  *int64    `json:"promo_price"`
	Address     *string   `json:"address"`
	Distance    *float32  `json:"distance"`
	Description *string   `json:"description"`
	Latitude    *string   `json:"latitude"`
	Longitude   *string   `json:"longitude"`
	Rating      *float32  `json:"rating"`
}
