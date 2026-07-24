package pkg

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(stc any) (map[string]string, error) {
	if err := validate.Struct(stc); err != nil {
		return formatValidationError(err)
	}
	return nil, nil
}

func formatValidationError(err error) (map[string]string, error) {
	var validationError validator.ValidationErrors
	if !errors.As(err, &validationError) {
		return nil, err
	}
	messages := make(map[string]string, len(validationError))
	for _, field := range validationError {
		messages[toSnakeCase(field.Field())] = formatMessage(field)
	}
	return messages, nil
}

func formatMessage(field validator.FieldError) string {
	switch field.Tag() {
	case "required":
		return "is required"
	case "min":
		return "must be at least " + field.Param() + " characters"
	case "max":
		return "must be at most " + field.Param() + " characters"
	default:
		return "is invalid"
	}
}
