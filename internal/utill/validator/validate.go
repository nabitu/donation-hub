package validator

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate() *validator.Validate {
	return validate
}
