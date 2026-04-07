package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns user-friendly error
func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Parse validation errors
	validationErrors := err.(validator.ValidationErrors)

	var messages []string
	for _, e := range validationErrors {
		messages = append(messages, formatValidationError(e))
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

func formatValidationError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
