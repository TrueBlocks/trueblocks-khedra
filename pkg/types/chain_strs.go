package types

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// CleanChainString processes and ensures the correctness of a chain string.
//   - Uses splitChainString to validate and clean the input.
//   - Guarantees that "mainnet" appears at the front of the returned `chains` string,
//     appending it if not already included.
//   - The `chains` string includes all valid, deduplicated chains starting with "mainnet".
//   - The `targets` string preserves the validated input order of chains.
func CleanChainString(input string) (string, string, error) {
	targets, err := splitChainString(input)
	if err != nil {
		return "", "", fmt.Errorf("invalid chain string: %v", err)
	}

	chainStrs := []string{"mainnet"}

	for _, chain := range targets {
		if chain != "mainnet" {
			chainStrs = append(chainStrs, chain)
		}
	}

	return strings.Join(chainStrs, ","), strings.Join(targets, ","), nil
}

// splitChainString validates and processes a comma-separated string of chains.
// - Trims leading/trailing whitespace from each chain.
// - Ensures no internal whitespace within each chain.
// - Validates that each chain contains only alphanumeric characters, dashes, and underscores.
// - Removes duplicates while preserving order.
// Returns a slice of valid chains or an error if validation fails.
func splitChainString(input string) ([]string, error) {
	validChainRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	parts := strings.Split(input, ",")
	var cleanedParts []string
	haveMap := map[string]bool{}

	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if len(trimmedPart) == 0 {
			continue
		}
		for _, r := range trimmedPart {
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("%w: '%s'", ErrInternalWhitespace, trimmedPart)
			}
		}
		if !validChainRegex.MatchString(trimmedPart) {
			return nil, fmt.Errorf("%w: '%s'", ErrInvalidCharacter, trimmedPart)
		}
		if !haveMap[trimmedPart] {
			cleanedParts = append(cleanedParts, trimmedPart)
			haveMap[trimmedPart] = true
		}
	}

	if len(cleanedParts) == 0 {
		return nil, ErrEmptyResult
	}

	return cleanedParts, nil
}

var (
	ErrInternalWhitespace = fmt.Errorf("invalid chain string: internal whitespace in part")
	ErrInvalidCharacter   = fmt.Errorf("invalid chain string: invalid character in part")
	ErrEmptyResult        = fmt.Errorf("invalid chain string: no valid chains found")
)
