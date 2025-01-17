package validate

import "fmt"

type validationError struct {
	errors []error
}

func (ve *validationError) Error() string {
	var result string
	for i, err := range ve.errors {
		result += fmt.Sprintf("Error %d: %s\n", i+1, err.Error())
	}
	return result
}

func NewValidationError(errors []error) *validationError {
	return &validationError{errors: errors}
}
