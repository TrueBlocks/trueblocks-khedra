package validate

import (
	"fmt"
	"reflect"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

func Validate(input interface{}) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	}

	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("Validate requires a pointer, got %T", input)
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("Validate requires a pointer to a struct, got pointer to %T", val.Interface())
	}

	var fieldValidators []FieldValidator
	collectFieldValidators(val, "", &fieldValidators, input)

	var errs []error
	for _, fv := range fieldValidators {
		for _, directive := range fv.directives {
			fv.validatorName, fv.tagArg = splitDirective(directive)
			fn, ok := validatorRegistry[fv.validatorName]
			if !ok {
				fmt.Println(colors.Red, "unknown validator", fv.validatorName, colors.Off)
				continue
			}
			if err := fn(fv); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}
