package types

import (
	"os"
	"reflect"
	"strings"
)

type KeyType string

const (
	InStruct KeyType = "in_struct"
	InEnv    KeyType = "in_env"
)

func GetEnvironmentKeys(inputStruct interface{}, kt KeyType) []string {
	return getEnvironmentKeys(inputStruct, kt)
}

func getEnvironmentKeys(inputStruct interface{}, kt KeyType) []string {
	rawKeys := getEnvironmentKeysInternal(inputStruct)
	var transformedKeys []string

	skip := func(key string) bool {
		filters := []string{
			"_NAME",
			"_API_BATCHSIZE",
			"_API_SLEEP",
			"_IPFS_BATCHSIZE",
			"_IPFS_SLEEP",
			"_MONITOR_PORT",
			"_SCRAPER_PORT",
		}
		for _, filter := range filters {
			if strings.HasSuffix(key, filter) {
				return true
			}
		}
		return false
	}

	for _, key := range rawKeys {
		transformedKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		transformedKey = strings.ReplaceAll(transformedKey, ",OMITEMPTY", "")
		transformedKey = "TB_KHEDRA_" + transformedKey

		if kt == InEnv {
			if _, exists := os.LookupEnv(transformedKey); !exists {
				continue
			}
		}

		if !skip(transformedKey) {
			transformedKeys = append(transformedKeys, transformedKey)
		}
	}

	return transformedKeys
}

func getEnvironmentKeysInternal(inputStruct interface{}) []string {
	var keys []string

	extractKeys := func(val reflect.Value, parentKey string) {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return
		}

		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldVal := val.Field(i)

			if tag, ok := field.Tag.Lookup("koanf"); ok {
				currentKey := tag
				if parentKey != "" {
					currentKey = parentKey + "." + tag
				}

				if fieldVal.Kind() == reflect.Map {
					for _, mapKey := range fieldVal.MapKeys() {
						mapValue := fieldVal.MapIndex(mapKey)
						if mapValue.Kind() == reflect.Ptr && !mapValue.IsNil() {
							mapValue = mapValue.Elem()
						}
						if mapValue.Kind() == reflect.Struct {
							subKeys := getEnvironmentKeysInternal(mapValue.Interface())
							for _, subKey := range subKeys {
								keys = append(keys, currentKey+"."+mapKey.String()+"."+subKey)
							}
						}
					}
				} else if fieldVal.Kind() == reflect.Struct || (fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().Kind() == reflect.Struct) {
					subKeys := getEnvironmentKeysInternal(fieldVal.Interface())
					for _, subKey := range subKeys {
						keys = append(keys, currentKey+"."+subKey)
					}
				} else {
					keys = append(keys, currentKey)
				}
			}
		}
	}

	val := reflect.ValueOf(inputStruct)
	extractKeys(val, "")
	return keys
}
