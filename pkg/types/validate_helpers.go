package types

import (
	"fmt"
	"net/url"
	"strings"
)

// Enum validates that a value is one of the allowed values
func Enum(value string, allowedValues []string) (bool, error) {
	for _, allowed := range allowedValues {
		if value == allowed {
			return true, nil
		}
	}
	return false, fmt.Errorf("value '%s' not in allowed values: %s", value, strings.Join(allowedValues, ", "))
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check that it has a scheme and host
	return u.Scheme != "" && u.Host != ""
}
