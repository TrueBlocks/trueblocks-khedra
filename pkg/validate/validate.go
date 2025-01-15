package validate

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

// FieldValidator holds everything needed to validate a single field.
type FieldValidator struct {
	fieldValue reflect.Value
	typeName   string
	fieldName  string
	context    string
	tagArg     string
	root       interface{}
	directives []string
}

func (fv FieldValidator) Root() interface{} {
	return fv.root
}

func (fv FieldValidator) Context() string {
	return fv.context
}

// ValidatorFunc defines the signature for validation functions.
type ValidatorFunc func(fv FieldValidator) error

// Validator registry
var validatorRegistry = make(map[string]ValidatorFunc)

func init() {
	validatorRegistry["oneof"] = oneofValidator
	validatorRegistry["min"] = minValidator
	validatorRegistry["max"] = maxValidator
	validatorRegistry["endswith"] = endswithValidator
	validatorRegistry["folder_exists"] = folderExistsValidator
	validatorRegistry["file_exists"] = fileExistsValidator
	validatorRegistry["strict_url"] = strictURLValidator
	validatorRegistry["is_writable"] = isWritableValidator
	validatorRegistry["opt_max"] = optMaxValidator
	validatorRegistry["opt_min"] = optMinValidator
	validatorRegistry["ping_one"] = pingOneValidator
	validatorRegistry["required"] = requiredValidator
	validatorRegistry["req_if_enabled"] = reqIfEnabledValidator
	validatorRegistry["dive"] = func(fv FieldValidator) error { return nil }
}

// RegisterValidator registers a new validator function if it does not
// already exist. If it does, it returns an error.
func RegisterValidator(name string, fn ValidatorFunc) error {
	if _, ok := validatorRegistry[name]; ok {
		return fmt.Errorf("validator %q already exists", name)
	}
	validatorRegistry[name] = fn
	return nil
}

// CollectFieldValidators recursively collects all fields with validation tags.
func collectFieldValidators(val reflect.Value, context string, fieldValidators *[]FieldValidator, root interface{}) {
	switch val.Kind() {
	case reflect.Ptr:
		if !val.IsNil() {
			collectFieldValidators(val.Elem(), context, fieldValidators, root)
		}
	case reflect.Struct:
		t := val.Type()
		typeName := t.Name()
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			fieldType := t.Field(i)
			if fieldType.PkgPath != "" {
				continue
			}
			validateTag := fieldType.Tag.Get("validate")
			if validateTag == "" {
				collectFieldValidators(fieldVal, context, fieldValidators, root)
				continue
			}
			directives := parseValidateTag(validateTag)
			preDive, postDive := splitDirectives(directives)
			if len(preDive) > 0 {
				fv := FieldValidator{
					fieldValue: fieldVal,
					typeName:   typeName,
					fieldName:  fieldType.Name,
					context:    context,
					directives: preDive,
					root:       root, // Attach the root object
				}
				*fieldValidators = append(*fieldValidators, fv)
			}
			if len(postDive) == 0 {
				collectFieldValidators(fieldVal, context, fieldValidators, root)
				continue
			}
			switch fieldVal.Kind() {
			case reflect.Slice, reflect.Array:
				if !fieldVal.IsValid() {
					continue
				}
				for j := 0; j < fieldVal.Len(); j++ {
					elemVal := fieldVal.Index(j)
					elemContext := fmt.Sprintf("%s[%d]", context, j)
					if elemVal.Kind() == reflect.Struct || (elemVal.Kind() == reflect.Ptr && !elemVal.IsNil()) {
						collectFieldValidators(elemVal, elemContext, fieldValidators, root)
						fv := FieldValidator{
							fieldValue: elemVal,
							typeName:   typeName,
							fieldName:  fieldType.Name,
							context:    elemContext,
							directives: postDive,
						}
						*fieldValidators = append(*fieldValidators, fv)
					} else {
						fv := FieldValidator{
							fieldValue: elemVal,
							typeName:   typeName,
							fieldName:  fieldType.Name,
							context:    elemContext,
							directives: postDive,
						}
						*fieldValidators = append(*fieldValidators, fv)
					}
				}
			case reflect.Map:
				for _, key := range fieldVal.MapKeys() {
					elemVal := fieldVal.MapIndex(key)
					elemContext := fmt.Sprintf("%s[%v]", context, key.Interface())
					if elemVal.Kind() == reflect.Struct || (elemVal.Kind() == reflect.Ptr && !elemVal.IsNil()) {
						collectFieldValidators(elemVal, elemContext, fieldValidators, root)
						fv := FieldValidator{
							fieldValue: elemVal,
							typeName:   typeName,
							fieldName:  fieldType.Name,
							context:    elemContext,
							directives: postDive,
						}
						*fieldValidators = append(*fieldValidators, fv)
					} else {
						fv := FieldValidator{
							fieldValue: elemVal,
							typeName:   typeName,
							fieldName:  fieldType.Name,
							context:    elemContext,
							directives: postDive,
						}
						*fieldValidators = append(*fieldValidators, fv)
					}
				}
			default:
				collectFieldValidators(fieldVal, context, fieldValidators, root)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			collectFieldValidators(val.Index(i), fmt.Sprintf("%s[%d]", context, i), fieldValidators, root)
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			collectFieldValidators(val.MapIndex(key), fmt.Sprintf("%s[%v]", context, key.Interface()), fieldValidators, root)
		}
	}
}

func splitDirectives(directives []string) (preDive, postDive []string) {
	for i, directive := range directives {
		if directive == "dive" {
			return directives[:i], directives[i+1:]
		}
	}
	return directives, nil
}

func parseValidateTag(tag string) []string {
	parts := strings.Split(tag, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func splitDirective(directive string) (string, string) {
	if idx := strings.Index(directive, "="); idx != -1 {
		return directive[:idx], directive[idx+1:]
	}
	return directive, ""
}

type ValidationError struct {
	errors []error
}

func (ve *ValidationError) Error() string {
	var result string
	for i, err := range ve.errors {
		result += fmt.Sprintf("Error %d: %s\n", i+1, err.Error())
	}
	return result
}

func NewValidationError(errors []error) *ValidationError {
	return &ValidationError{errors: errors}
}

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
			validatorName, tagArg := splitDirective(directive)
			fn, ok := validatorRegistry[validatorName]
			if !ok {
				continue
			}
			fv.tagArg = tagArg
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

// Validator implementations

func oneofValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "oneof", "is not a string", fv.fieldValue.Kind().String())
	}

	validValues := strings.Fields(fv.tagArg)
	for _, test := range validValues {
		if value == test {
			return Passed(fv, "oneof", value, fv.tagArg)
		}
	}

	return Failed(fv, "oneof", fmt.Sprintf("must be one of %v", validValues), fmt.Sprintf("%q", value))
}

func minValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "min", "is not an integer", fv.fieldValue.Kind().String())
	}

	test, err := strconv.Atoi(fv.tagArg)
	if err != nil {
		return Failed(fv, "min", "invalid tag argument", fv.tagArg)
	}
	if value < int64(test) {
		return Failed(fv, "min", fmt.Sprintf("must be >= %d", test), fmt.Sprintf("%d", value))
	}

	return Passed(fv, "min", fmt.Sprintf("%d", value), fmt.Sprintf("%d", test))
}

func maxValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "max", "is not an integer", fv.fieldValue.Kind().String())
	}

	maxValue, err := strconv.Atoi(fv.tagArg)
	if err != nil {
		return Failed(fv, "max", "tag argument is not an integer", fv.tagArg)
	}
	if value > int64(maxValue) {
		return Failed(fv, "max", fmt.Sprintf("must be <= %d", maxValue), fmt.Sprintf("%q", value))
	}

	return Passed(fv, "max", fmt.Sprintf("%d", value), fmt.Sprintf("%d", maxValue))
}

func endswithValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "endswith", "is not a string", fv.fieldValue.Kind().String())
	}

	if !strings.HasSuffix(value, fv.tagArg) {
		return Failed(fv, "endswith", fmt.Sprintf("must end with %q", fv.tagArg), fmt.Sprintf("%q", value))
	}

	return Passed(fv, "endswith", value, fv.tagArg)
}

func folderExistsValidator(fv FieldValidator) error {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return Passed(fv, "folder_exists", "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "folder_exists", "is not a string", fv.fieldValue.Kind().String())
	}

	value = utils.ResolvePath(value)
	info, err := os.Stat(value)
	if os.IsNotExist(err) {
		return Failed(fv, "folder_exists", "folder does not exist", fmt.Sprintf("%q", value))
	}
	if err != nil {
		return Failed(fv, "folder_exists", "error accessing path", fmt.Sprintf("%q: %v", value, err))
	}

	if !info.IsDir() {
		return Failed(fv, "folder_exists", "path exists but is not a folder", fmt.Sprintf("%q", value))
	}

	return Passed(fv, "folder_exists", value, "")
}

func requiredValidator(fv FieldValidator) error {
	if !fv.fieldValue.IsValid() {
		return Failed(fv, "required", "is required", "")
	}

	switch fv.fieldValue.Kind() {
	case reflect.Slice, reflect.Array:
		if fv.fieldValue.Len() == 0 {
			return Failed(fv, "required", "cannot be empty", "")
		}
	}

	if fv.fieldValue.IsZero() {
		return Failed(fv, "required", "is required", "")
	}

	return Passed(fv, "required", "", "")
}

func fileExistsValidator(fv FieldValidator) error {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return Passed(fv, "file_exists", "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "file_exists", "is not a string", fv.fieldValue.Kind().String())
	}

	value = utils.ResolvePath(value)
	info, err := os.Stat(value)
	if err != nil || info.IsDir() {
		return Failed(fv, "file_exists", "file does not exist", value)
	}

	return Passed(fv, "file_exists", value, "")
}

func isWritableValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "is_writable", "is not a string", fv.fieldValue.Kind().String())
	}

	file, err := os.OpenFile(value, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return Failed(fv, "is_writable", "path is not writable", value)
	}
	file.Close()

	return Passed(fv, "is_writable", value, "")
}

func optMaxValidator(fv FieldValidator) error {
	if !fv.fieldValue.IsValid() || fv.fieldValue.IsZero() {
		return Passed(fv, "opt_max", "unset-ok", "")
	}

	return maxValidator(fv)
}

func optMinValidator(fv FieldValidator) error {
	if !fv.fieldValue.IsValid() || fv.fieldValue.IsZero() {
		return Passed(fv, "opt_min", "unset-ok", "")
	}

	return minValidator(fv)
}

type Enabler interface {
	IsEnabled() bool
}

func reqIfEnabledValidator(fv FieldValidator) error {
	en, ok := fv.root.(Enabler)
	if !ok {
		// If the structure is not an Enabler (for example, General), assume it's enabled
	} else if !en.IsEnabled() {
		return Passed(fv, "req_if_enabled", "not-enabled", "")
	}

	return requiredValidator(fv)
}

func strictURLValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "strict_url", "is not a string", fv.fieldValue.Kind().String())
	}

	if _, err := url.ParseRequestURI(value); err != nil {
		return Failed(fv, "strict_url", "is not a valid URL", value)
	}

	return Passed(fv, "strict_url", value, "")
}

// Not implemented

func pingOneValidator(fv FieldValidator) error {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return Passed(fv, "ping_one", "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return Failed(fv, "ping_one", "is not a string", fv.fieldValue.Kind().String())
	}
	return Passed(fv, "ping_one", value, "")
}

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

func Passed(fv FieldValidator, validatorName, value, test string) error {
	// c := fmt.Sprintf(" context=%q", fv.context)
	// if fv.context == "" {
	// 	c = ""
	// }
	// fmt.Printf("%s%-20.20s [%-13.13s] PASSED (value=%q test=%q%s)%s\n", colors.Green, fv.typeName+"."+fv.fieldName, validatorName, value, test, c, colors.Off)
	return nil
}

func Failed(fv FieldValidator, validatorName, errStr, got string) error {
	c := fmt.Sprintf(" (context=%q)", fv.context)
	if fv.context == "" {
		c = ""
	}
	return fmt.Errorf("\n%s[%-13.13s] FAILED for %s.%s%s %s (got %s)%s", colors.Red, validatorName, fv.typeName, fv.fieldName, c, errStr, got, colors.Off)
}
