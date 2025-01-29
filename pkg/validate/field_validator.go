package validate

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldValidator holds everything needed to validate a single field.
type FieldValidator struct {
	fieldValue    reflect.Value
	validatorName string
	typeName      string
	fieldName     string
	context       string
	tagArg        string
	root          interface{}
	directives    []string
}

func NewFieldValidator(name, typ, fld, ctx string) FieldValidator {
	return FieldValidator{
		validatorName: name,
		typeName:      typ,
		fieldName:     fld,
		context:       ctx,
	}
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
var ValidatorRegistry = make(map[string]ValidatorFunc)

func init() {
	ValidatorRegistry["oneof"] = oneofValidator
	ValidatorRegistry["min"] = minValidator
	ValidatorRegistry["max"] = maxValidator
	ValidatorRegistry["endswith"] = endswithValidator
	ValidatorRegistry["folder_exists"] = folderExistsValidator
	ValidatorRegistry["strict_url"] = strictURLValidator
	ValidatorRegistry["non_zero"] = nonZeroValidator
	ValidatorRegistry["required"] = requiredValidator
	ValidatorRegistry["req_if_enabled"] = reqIfEnabledValidator
	ValidatorRegistry["dive"] = func(fv FieldValidator) error { return nil }
}

// RegisterValidator registers a new validator function if it does not
// already exist. If it does, it returns an error.
func RegisterValidator(name string, fn ValidatorFunc) error {
	if _, ok := ValidatorRegistry[name]; ok {
		return fmt.Errorf("validator %q already exists", name)
	}
	ValidatorRegistry[name] = fn
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
