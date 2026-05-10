package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// InitValidator is the hook for registering custom validation rules.
func InitValidator() {
}

// HandleValidatorError converts validation errors into a concise client message.
func HandleValidatorError(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "invalid request parameters"
	}

	err := errs[0]
	field := strings.ToLower(err.Field())
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
