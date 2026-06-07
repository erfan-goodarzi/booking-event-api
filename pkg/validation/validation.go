package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func FormatValidationErrors(err error) map[string]string {
	fields := make(map[string]string)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			fields[e.Field()] = e.Tag()
		}
	}
	return fields
}
