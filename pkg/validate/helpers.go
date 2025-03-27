package validate

import (
	"fmt"
	"reflect"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

// Helpers

func getStringValue(fieldVal reflect.Value) (string, error) {
	if fieldVal.Kind() != reflect.String {
		return "", fmt.Errorf("placeholder")
	}
	return fieldVal.String(), nil
}

func getIntValue(fieldVal reflect.Value) (int64, error) {
	if fieldVal.Kind() != reflect.Int {
		return 0, fmt.Errorf("placeholder")
	}
	return fieldVal.Int(), nil
}

func Passed(fv FieldValidator, value, test string) error {
	_ = fv
	_ = value
	_ = test
	// c := fmt.Sprintf(" context=%q", fv.context)
	// if fv.context == "" {
	// 	c = ""
	// }
	// fmt.Printf("%s%-20.20s [%-13.13s] PASSED (value=%q test=%q%s)%s\n", colors.Green, fv.typeName+"."+fv.fieldName, fv.validatorName, value, test, c, colors.Off)
	return nil
}

func Failed(fv FieldValidator, errStr, got string) error {
	c := fmt.Sprintf(" (context=%q)", fv.context)
	if fv.context == "" {
		c = ""
	}
	return fmt.Errorf("\n%s[%-13.13s] FAILED for %s.%s%s %s (got %s)%s", colors.Red, fv.validatorName, fv.typeName, fv.fieldName, c, errStr, got, colors.Off)
}
