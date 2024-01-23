package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ceo-suite/global/helper"
	"github.com/ceo-suite/internal/user"
)

type forgotPasswordHandler struct {
	user user.Service
}

type forgotPasswordRequesData struct {
	Email string `json:"email"`
}

func (h *forgotPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle based on HTTP request method
	switch r.Method {
	case http.MethodPost:
		h.handleForgotPassword(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *forgotPasswordHandler) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		source     string          // stores request source
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[User HTTP][handleForgotPassword] Failed to forgot password. Source: %s, Err: %s\n", source, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		// read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// unmarshall body
		var data forgotPasswordRequesData
		err = json.Unmarshal(body, &data)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// Forgot password
		err = h.user.ForgotPassword(ctx, data.Email)
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
				log.Printf("[User HTTP][handleForgotPassword] Internal error from ForgotPassword. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- data.Email
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case userID := <-resChan:
		res := helper.ResponseEnvelope{
			Data: userID,
		}
		resBody, err = json.Marshal(res)
	}
}
