package user

import "errors"

// Followings are the known errors returned from user.
var (
	// ErrUserAlreadyExist is returned when the given
	// user already exist based on the predefined
	// unique constraints.
	ErrUserAlreadyExist = errors.New("user already exist")

	// ErrDataNotFound is returned when the desired data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrExpiredToken is returned when the given token is
	// expired.
	ErrExpiredToken = errors.New("expired token")

	// ErrInvalidPassword is returned when the given password
	// is invalid.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrInvalidToken is returned when the given token is
	// invalid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidEmail is returned when the given email is
	// invalid.
	ErrInvalidEmail = errors.New("invalid email")

	// ErrInvalidPICName is returned when the given PIC name is
	// invalid.
	ErrInvalidPICName = errors.New("invalid PIC name")

	// ErrInvalidCompanyName is returned when the given company name is
	// invalid.
	ErrInvalidCompanyName = errors.New("invalid company name")

	// ErrInvalidUserID is returned when the given user ID is
	// invalid.
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrInvalidPhoneNumber is returned when the given phone number is
	// invalid.
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)
