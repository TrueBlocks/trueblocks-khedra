package types

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// oneofValidator ensures the field value matches one of the allowed values specified in the tag argument.
func oneofValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not a string", fv.fieldValue.Kind().String())
	}

	validValues := strings.Fields(fv.tagArg)
	for _, test := range validValues {
		if value == test {
			return Passed(fv, value, fv.tagArg)
		}
	}

	return Failed(fv, fmt.Sprintf("must be one of %v", validValues), fmt.Sprintf("%q", value))
}

// minValidator ensures the field value is an integer greater than or equal to the minimum specified in the tag argument.
func minValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not an integer", fv.fieldValue.Kind().String())
	}

	test, err := strconv.Atoi(fv.tagArg)
	if err != nil {
		return Failed(fv, "invalid tag argument", fv.tagArg)
	}
	if value < int64(test) {
		return Failed(fv, fmt.Sprintf("must be >= %d", test), fmt.Sprintf("%d", value))
	}

	return Passed(fv, fmt.Sprintf("%d", value), fmt.Sprintf("%d", test))
}

// maxValidator ensures the field value is an integer less than or equal to the maximum specified in the tag argument.
func maxValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not an integer", fv.fieldValue.Kind().String())
	}

	maxValue, err := strconv.Atoi(fv.tagArg)
	if err != nil {
		return Failed(fv, "tag argument is not an integer", fv.tagArg)
	}
	if value > int64(maxValue) {
		return Failed(fv, fmt.Sprintf("must be <= %d", maxValue), fmt.Sprintf("%q", value))
	}

	return Passed(fv, fmt.Sprintf("%d", value), fmt.Sprintf("%d", maxValue))
}

// endswithValidator ensures the field value is a string that ends with the suffix specified in the tag argument.
func endswithValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not a string", fv.fieldValue.Kind().String())
	}

	if !strings.HasSuffix(value, fv.tagArg) {
		return Failed(fv, fmt.Sprintf("must end with %q", fv.tagArg), fmt.Sprintf("%q", value))
	}

	return Passed(fv, value, fv.tagArg)
}

// folderExistsValidator ensures the field value is a string representing a path to an existing
// folder or creates a writable folder, if not.
func folderExistsValidator(fv FieldValidator) error {
	isTesting := base.IsTestMode()
	if isTesting {
		return Passed(fv, "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not a string", fv.fieldValue.Kind().String())
	}

	sanitizedValue := utils.ToValidPath(value)
	if sanitizedValue != value {
		return Failed(fv, "invalid characters in path", fmt.Sprintf("%q-%q", value, sanitizedValue))
	}

	value = utils.ResolvePath(value)
	if info, err := os.Stat(value); err == nil && !info.IsDir() {
		return Failed(fv, "path exists but is not a folder", fmt.Sprintf("%q", value))
	}

	if !coreFile.FolderExists(value) {
		err := coreFile.EstablishFolder(value)
		if err != nil {
			return Failed(fv, err.Error(), fmt.Sprintf("%q", value))
		}
	}

	_, err = os.Stat(value)
	if err != nil {
		return Failed(fv, "error accessing path", fmt.Sprintf("%q: %v", value, err))
	}

	return Passed(fv, value, "")
}

// requiredValidator ensures the field value is non-zero, non-empty, and valid based on its type.
func requiredValidator(fv FieldValidator) error {
	if !fv.fieldValue.IsValid() {
		return Failed(fv, "is required", "")
	}

	switch fv.fieldValue.Kind() {
	case reflect.Slice, reflect.Array:
		if fv.fieldValue.Len() == 0 {
			return Failed(fv, "cannot be empty", "")
		}
	}

	if fv.fieldValue.IsZero() {
		return Failed(fv, "is required", "")
	}

	return Passed(fv, "", "")
}

type Enabler interface {
	IsEnabled() bool
}

// reqIfEnabledValidator ensures the field value is required if the root structure implements
// Enabler and is enabled. If it's not an Enabler, it assumes the structure passes.
func reqIfEnabledValidator(fv FieldValidator) error {
	en, ok := fv.root.(Enabler)
	if !ok {
		// If the structure is not an Enabler (for example, General), assume it's enabled
	} else if !en.IsEnabled() {
		return Passed(fv, "not-enabled", "")
	}

	return requiredValidator(fv)
}

// strictURLValidator ensures the field value is a string representing a valid URL.
func strictURLValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not a string", fv.fieldValue.Kind().String())
	}

	if _, err := url.ParseRequestURI(value); err != nil {
		return Failed(fv, "is not a valid URL", value)
	}

	return Passed(fv, value, "")
}

// nonZeroValidator ensures the field value is an integer greater than zero.
func nonZeroValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is not an integer", fv.fieldValue.Kind().String())
	}

	if value == 0 {
		return Failed(fv, "must be non-zero", fmt.Sprintf("%d", value))
	}

	return Passed(fv, fmt.Sprintf("%d", value), "")
}
