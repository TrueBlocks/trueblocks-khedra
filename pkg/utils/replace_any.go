package utils

import "strings"

func RemoveAny(A, B string) string {
	result := strings.Builder{}
	toRemove := make(map[rune]struct{})
	for _, char := range B {
		toRemove[char] = struct{}{}
	}
	for _, char := range A {
		if _, exists := toRemove[char]; !exists {
			result.WriteRune(char)
		}
	}
	return result.String()
}
