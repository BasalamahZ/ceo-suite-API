package service

import (
	"context"
	"net/mail"

	emailClient "github.com/ceo-suite/global/email"
	"github.com/ceo-suite/internal/user"
)

// CreateUser create new user as given.
// CreateUser expects the given  user ID in the given
// user already assigned.
func (s *service) CreateUser(ctx context.Context, reqUser user.User) (int64, error) {
	// validate fields
	err := validateUser(reqUser)
	if err != nil {
		return 0, err
	}

	// hash password
	hash, err := generateHash(reqUser.Password, s.config.PasswordSalt)
	if err != nil {
		return 0, err
	}

	// modify fields
	reqUser.Password = hash
	reqUser.CreateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return 0, err
	}

	// inserts user in pgstore
	userID, err := pgStoreClient.CreateUser(ctx, reqUser)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// UpdatePassword, for the given user ID, updates user's
// password with the new password. Before updating, it
// checks whether the current password are correct.
func (s *service) UpdatePassword(ctx context.Context, userID int64, newPassword string, currentPassword string) error {
	if newPassword == "" || currentPassword == "" {
		return user.ErrInvalidPassword
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// get user current data
	current, err := pgStoreClient.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// check current password
	err = s.checkPassword(ctx, current, currentPassword)
	if err != nil {
		return err
	}

	// hash password
	hash, err := generateHash(newPassword, s.config.PasswordSalt)
	if err != nil {
		return err
	}

	// update fields
	current.Password = hash
	current.UpdateTime = s.timeNow()

	// update user
	err = pgStoreClient.UpdateUser(ctx, current)
	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword, for the given email to verification email
// and then can reset user's password
func (s *service) ForgotPassword(ctx context.Context, email string) error {
	// validate the given values
	if email == "" {
		return user.ErrInvalidEmail
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// check user
	_, err = pgStoreClient.GetUserByEmail(ctx, email)
	if err != nil {
		return user.ErrDataNotFound
	}

	// go sendEmail(email)
	mail := emailClient.NewMailClient()
	mail.SetSender("basalamah76@gmail.com")
	mail.SetReciever(email)
	mail.SetSubject("Forgot Password")
	if err := mail.SendMail(); err != nil {
		return err
	}

	return nil
}

// ResetPassword, for the given user ID, updates user's
// password with the new password without checking the
// current password.
func (s *service) ResetPassword(ctx context.Context, userID int64, newPassword string) error {
	// validate the given values
	if userID <= 0 {
		return user.ErrInvalidUserID
	}
	if newPassword == "" {
		return user.ErrInvalidPassword
	}

	// hash password
	hash, err := generateHash(newPassword, s.config.PasswordSalt)
	if err != nil {
		return err
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// get user current data
	current, err := pgStoreClient.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// update fields
	current.Password = hash
	current.UpdateTime = s.timeNow()

	// update user
	err = pgStoreClient.UpdateUser(ctx, current)
	if err != nil {
		return err
	}

	return nil
}

// checkPassword checks the given password with the password
// stored in store for the given user ID.
func (s *service) checkPassword(ctx context.Context, data user.User, password string) error {
	// compare paswords
	match, err := compareHash(password, data.Password, s.config.PasswordSalt)
	if err != nil {
		return err
	}

	// wrong password
	if !match {
		return user.ErrInvalidPassword
	}

	return nil
}

// send email to the given email for verification
func sendEmail(email string) error {
	mail := emailClient.NewMailClient()
	mail.SetSender("basalamah76@gmail.com")
	mail.SetReciever(email)
	mail.SetSubject("Forgot Password")
	if err := mail.SendMail(); err != nil {
		return err
	}

	return nil
}

// validateUser validates fields of the given user
// whether its comply the predetermined rules.
func validateUser(reqUser user.User) error {
	_, err := mail.ParseAddress(reqUser.Email)
	if reqUser.Email == "" || err != nil {
		return user.ErrInvalidEmail
	}

	if reqUser.PICName == "" {
		return user.ErrInvalidPICName
	}

	if reqUser.CompanyName == "" {
		return user.ErrInvalidCompanyName
	}

	if reqUser.PhoneNumber == "" || len(reqUser.PhoneNumber) < 10 || len(reqUser.PhoneNumber) > 13 {
		return user.ErrInvalidPhoneNumber
	}

	return nil
}
