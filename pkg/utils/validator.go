package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// InitValidator initializes the validator instance
func InitValidator() {
	validate = validator.New()
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) map[string]string {
	v := GetValidator()
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = formatValidationError(err)
	}
	return errors
}

// formatValidationError formats a validation error into a user-friendly message
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	case "oneof":
		return err.Field() + " must be one of: " + err.Param()
	default:
		return err.Field() + " is invalid"
	}
}
