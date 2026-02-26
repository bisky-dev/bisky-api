package auth

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var authValidator = validator.New()

const (
	emailValidationRule    = "required,email,max=254"
	passwordValidationRule = "required,min=8,max=128"
)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateRegisterRequest(req registerRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	return validatePassword(req.Password)
}

func validateLoginRequest(req loginRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	return validatePassword(req.Password)
}

func validateForgotPasswordRequest(req forgotPasswordRequest) error {
	return validateEmail(req.Email)
}

func validateEmail(email string) error {
	return validateVar(email, emailValidationRule, "email is invalid")
}

func validatePassword(password string) error {
	return validateVar(password, passwordValidationRule, "password is invalid")
}

func validateVar(value string, rule, message string) error {
	if err := authValidator.Var(value, rule); err != nil {
		return errors.New(message)
	}
	return nil
}
