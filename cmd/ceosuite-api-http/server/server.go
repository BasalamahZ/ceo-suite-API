package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ceo-suite/cmd/ceosuite-api-http/config"
	"github.com/ceo-suite/internal/booking"
	bookinghttphandler "github.com/ceo-suite/internal/booking/handler/http"
	bookingservice "github.com/ceo-suite/internal/booking/service"
	bookingpgstore "github.com/ceo-suite/internal/booking/store/postgresql"
	"github.com/ceo-suite/internal/product"
	producthttphandler "github.com/ceo-suite/internal/product/handler/http"
	productservice "github.com/ceo-suite/internal/product/service"
	productpgstore "github.com/ceo-suite/internal/product/store/postgresql"
	"github.com/ceo-suite/internal/user"
	userhttphandler "github.com/ceo-suite/internal/user/handler/http"
	userservice "github.com/ceo-suite/internal/user/service"
	userpgstore "github.com/ceo-suite/internal/user/store/postgresql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Following constants are the possible exit code returned
// when running a server.
const (
	CodeSuccess = iota
	CodeBadConfig
	CodeFailServeHTTP
)

// Run creates a server and starts the server.
//
// Run returns a status code suitable for os.Exit() argument.
func Run() int {
	s, err := new()
	if err != nil {
		return CodeBadConfig
	}

	return s.start()
}

// server is the long-runnning application.
type server struct {
	srv      *http.Server
	handlers []handler
}

// handler provides mechanism to start HTTP handler. All HTTP
// handlers must implements this interface.
type handler interface {
	Start(multiplexer *mux.Router) error
}

// new creates and returns a new server.
func new() (*server, error) {
	s := &server{
		srv: &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// connect to database
	db, err := sqlx.Connect("postgres", config.BaseConfig())
	if err != nil {
		log.Printf("[ceo-suite-api-http] failed to connect database: %s\n", err.Error())
		return nil, fmt.Errorf("failed to connect database: %s", err.Error())
	}

	// initialize user service
	var userSvc user.Service
	{
		pgStore, err := userpgstore.New(db)
		if err != nil {
			log.Printf("[user-api-http] failed to initialize user postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize user postgresql store: %s", err.Error())
		}

		svcOptions := []userservice.Option{}
		svcOptions = append(svcOptions, userservice.WithConfig(userservice.Config{
			PasswordSalt:   os.Getenv("PasswordSalt"),
			TokenSecretKey: os.Getenv("TokenSecretKey"),
		}))

		userSvc, err = userservice.New(pgStore, svcOptions...)
		if err != nil {
			log.Printf("[user-api-http] failed to initialize user service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize user service: %s", err.Error())
		}
	}

	// initialize product service
	var productSvc product.Service
	{
		pgStore, err := productpgstore.New(db)
		if err != nil {
			log.Printf("[product-api-http] failed to initialize product postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize product postgresql store: %s", err.Error())
		}

		productSvc, err = productservice.New(pgStore)
		if err != nil {
			log.Printf("[tenant-api-http] failed to initialize product service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize product service: %s", err.Error())
		}
	}

	// initialize booking service
	var bookingSvc booking.Service
	{
		pgStore, err := bookingpgstore.New(db)
		if err != nil {
			log.Printf("[booking-api-http] failed to initialize booking postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize booking postgresql store: %s", err.Error())
		}

		bookingSvc, err = bookingservice.New(pgStore)
		if err != nil {
			log.Printf("[tenant-api-http] failed to initialize booking service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize booking service: %s", err.Error())
		}
	}

	// initialize user HTTP handler
	{
		identities := []userhttphandler.HandlerIdentity{
			userhttphandler.HandlerPassword,
			userhttphandler.HandlerToken,
			userhttphandler.HandlerLogin,
			userhttphandler.HandlerForgotPassword,
			userhttphandler.HandlerUsers,
		}

		userHTTP, err := userhttphandler.New(userSvc, identities)
		if err != nil {
			log.Printf("[user-api-http] failed to initialize user http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize user http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, userHTTP)
	}

	// initialize product HTTP handler
	{
		identities := []producthttphandler.HandlerIdentity{
			producthttphandler.HandlerProduct,
			producthttphandler.HandlerProducts,
		}

		productHTTP, err := producthttphandler.New(productSvc, userSvc, identities)
		if err != nil {
			log.Printf("[product-api-http] failed to initialize product http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize product http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, productHTTP)
	}

	// initialize booking HTTP handler
	{
		identities := []bookinghttphandler.HandlerIdentity{
			bookinghttphandler.HandlerBooking,
			bookinghttphandler.HandlerBookings,
		}

		bookingHTTP, err := bookinghttphandler.New(bookingSvc, userSvc, identities)
		if err != nil {
			log.Printf("[booking-api-http] failed to initialize booking http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize booking http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, bookingHTTP)
	}

	return s, nil
}

// start starts the given server.
func (s *server) start() int {
	log.Println("[ceo-suite-api-http] starting server...")

	// create multiplexer object
	rootMux := mux.NewRouter()
	appMux := rootMux.PathPrefix("/api").Subrouter()

	// starts handlers
	for _, h := range s.handlers {
		if err := h.Start(appMux); err != nil {
			log.Printf("[ceo-suite-api-http] failed to start handler: %s\n", err.Error())
			return CodeFailServeHTTP
		}
	}

	// endpoint checker
	appMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world! @ceo-suite")
	})

	// use middlewares to app mux only
	appMux.Use(corsMiddleware)

	// listen and serve
	log.Printf("[ceo-suite-api-http] Server is running at %s:%s", os.Getenv("ADDRESS"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("ADDRESS"), os.Getenv("PORT")), rootMux))

	return CodeSuccess
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
