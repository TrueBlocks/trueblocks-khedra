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
	validatorRegistry["dive"] = diveValidator
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

func Validate4(input interface{}) error {
	errs := Validate2(input)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func Validate2(input interface{}) []error {
	if input == nil {
		return []error{fmt.Errorf("input is nil")}
	}

	val := reflect.ValueOf(input)

	// Ensure we have a pointer
	if val.Kind() != reflect.Ptr {
		return []error{fmt.Errorf("Validate2 requires a pointer, got %T", input)}
	}

	// Dereference the pointer
	val = val.Elem()

	// Check that the dereferenced pointer is a struct
	if val.Kind() != reflect.Struct {
		return []error{fmt.Errorf("Validate2 requires a pointer to a struct, got pointer to %T", val.Interface())}
	}

	// Now proceed with collecting validators using 'val' and passing 'input' as the root.
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

	return errs
}

// Validator implementations

func oneofValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "oneof", "is not a string", fv.fieldValue.Kind().String())
	}

	tests := strings.Fields(fv.tagArg)
	for _, test := range tests {
		if value == test {
			return passed(fv, "oneof", value, fv.tagArg)
		}
	}

	return failed(fv, "oneof", fmt.Sprintf("must be one of %v", tests), fmt.Sprintf("%q", value))
}

func minValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "min", "is not an integer", fv.fieldValue.Kind().String())
	}

	test, err := strconv.Atoi(fv.tagArg)
	if err != nil {
		return failed(fv, "min", "invalid tag argument", fv.tagArg)
	}
	if value < int64(test) {
		return failed(fv, "min", fmt.Sprintf("must be >= %d", test), fmt.Sprintf("%d", value))
	}

	return passed(fv, "min", fmt.Sprintf("%d", value), fmt.Sprintf("%d", test))
}

func maxValidator(fv FieldValidator) error {
	value, err := getIntValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "max", "is not an integer", fv.fieldValue.Kind().String())
	}

	maxValue, err := strconv.Atoi(fv.tagArg)
	if err != nil {
		return failed(fv, "max", "tag argument is not an integer", fv.tagArg)
	}
	if value > int64(maxValue) {
		return failed(fv, "max", fmt.Sprintf("must be <= %d", maxValue), fmt.Sprintf("%q", value))
	}

	return passed(fv, "max", fmt.Sprintf("%d", value), fmt.Sprintf("%d", maxValue))
}

func endswithValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "endswith", "is not a string", fv.fieldValue.Kind().String())
	}

	if !strings.HasSuffix(value, fv.tagArg) {
		return failed(fv, "endswith", fmt.Sprintf("must end with %q", fv.tagArg), fmt.Sprintf("%q", value))
	}

	return passed(fv, "endswith", value, fv.tagArg)
}

func folderExistsValidator(fv FieldValidator) error {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return passed(fv, "ping_one", "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "folder_exists", "is not a string", fv.fieldValue.Kind().String())
	}

	value = utils.ResolvePath(value)
	if _, err := os.Stat(value); os.IsNotExist(err) {
		return failed(fv, "folder_exists", "folder does not exist", fmt.Sprintf("%q", value))
	}

	return passed(fv, "folder_exists", value, "")
}

func requiredValidator(fv FieldValidator) error {
	// Check if the field value is invalid
	if !fv.fieldValue.IsValid() {
		return failed(fv, "required", "is required", "")
	}

	// Handle specific kinds like slices, arrays
	switch fv.fieldValue.Kind() {
	case reflect.Slice, reflect.Array:
		if fv.fieldValue.Len() == 0 {
			return failed(fv, "required", "cannot be empty", "")
		}
	}

	// General zero-value check for other kinds
	if fv.fieldValue.IsZero() {
		return failed(fv, "required", "is required", "")
	}

	// fmt.Println(colors.Blue, fv.fieldName, fv.fieldValue, colors.Off)
	return passed(fv, "required", "", "")
}

func fileExistsValidator(fv FieldValidator) error {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return passed(fv, "ping_one", "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "file_exists", "is not a string", fv.fieldValue.Kind().String())
	}

	value = utils.ResolvePath(value)
	info, err := os.Stat(value)
	if err != nil || info.IsDir() {
		return failed(fv, "file_exists", "file does not exist", value)
	}

	return passed(fv, "file_exists", value, "")
}

func isWritableValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "is_writable", "is not a string", fv.fieldValue.Kind().String())
	}

	file, err := os.OpenFile(value, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return failed(fv, "is_writable", "path is not writable", value)
	}
	file.Close()

	return passed(fv, "is_writable", value, "")
}

func optMaxValidator(fv FieldValidator) error {
	if !fv.fieldValue.IsValid() || fv.fieldValue.IsZero() {
		return passed(fv, "opt_max", "unset-ok", "")
	}

	return maxValidator(fv)
}

func optMinValidator(fv FieldValidator) error {
	if !fv.fieldValue.IsValid() || fv.fieldValue.IsZero() {
		return passed(fv, "opt_min", "unset-ok", "")
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
		return passed(fv, "req_if_enabled", "not-enabled", "")
	}

	return requiredValidator(fv)
}

func strictURLValidator(fv FieldValidator) error {
	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "strict_url", "is not a string", fv.fieldValue.Kind().String())
	}

	if _, err := url.ParseRequestURI(value); err != nil {
		return failed(fv, "strict_url", "is not a valid URL", value)
	}

	return passed(fv, "strict_url", value, "")
}

// Not implemented

func pingOneValidator(fv FieldValidator) error {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return passed(fv, "ping_one", "skipped", "")
	}

	value, err := getStringValue(fv.fieldValue)
	if err != nil {
		return failed(fv, "ping_one", "is not a string", fv.fieldValue.Kind().String())
	}
	return passed(fv, "ping_one", value, "")
}

func diveValidator(fv FieldValidator) error {
	return nil
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
	return passed(fv, validatorName, value, test)
}

func passed(fv FieldValidator, validatorName, value, test string) error {
	// c := fmt.Sprintf(" context=%q", fv.context)
	// if fv.context == "" {
	// 	c = ""
	// }
	// fmt.Printf("%s%-20.20s [%-13.13s] PASSED (value=%q test=%q%s)%s\n", colors.Green, fv.typeName+"."+fv.fieldName, validatorName, value, test, c, colors.Off)
	return nil
}

func Failed(fv FieldValidator, validatorName, value, test string) error {
	return failed(fv, validatorName, value, test)
}

func failed(fv FieldValidator, validatorName, errStr, got string) error {
	c := fmt.Sprintf(" (context=%q)", fv.context)
	if fv.context == "" {
		c = ""
	}
	return fmt.Errorf("\n%s[%-13.13s] FAILED for %s.%s%s %s (got %s)%s", colors.Red, validatorName, fv.typeName, fv.fieldName, c, errStr, got, colors.Off)
}
