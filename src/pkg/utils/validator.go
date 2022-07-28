package utils

import "github.com/go-playground/validator/v10"

var myValidator = validator.New()

func GetValidator() *validator.Validate {
	return myValidator
}

func ValidatorErrors(err error) map[string]string {
	// Define fields map
	fields := map[string]string{}

	// Make error message for each invalid field
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
