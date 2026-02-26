package httpx

import "errors"

func ValidateVar(value any, rule, message string) error {
	if err := validate.Var(value, rule); err != nil {
		return errors.New(message)
	}
	return nil
}

func ValidateOptionalDate(value *string, message string) error {
	if value == nil {
		return nil
	}
	return ValidateVar(*value, "datetime=2006-01-02", message)
}
