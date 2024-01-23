package http

import (
	"context"
	"log"

	"github.com/ceo-suite/internal/user"
)

// checkAccessToken checks the given access token whether it
// is valid or not.
func checkAccessToken(ctx context.Context, svc user.Service, token, name string) error {
	_, err := svc.ValidateToken(ctx, token)
	if err != nil {
		log.Printf("[Booking HTTP][%s] Unauthorized error from ValidateToken. Err: %s\n", name, err.Error())
		return errUnauthorizedAccess
	}

	// if userID != tokenData.UserID {
	// 	return errInvalidUserID
	// }

	return nil
}
