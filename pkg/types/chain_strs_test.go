package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests derived from ai/TestDesign_chain_strs.go.md

func TestCleanChainString_HappyPaths(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantChains  string
		wantTargets string
	}{
		{name: "Add mainnet if missing", input: "goerli", wantChains: "mainnet,goerli", wantTargets: "goerli"},
		{name: "Preserve mainnet first", input: "mainnet,goerli", wantChains: "mainnet,goerli", wantTargets: "mainnet,goerli"},
		{name: "Deduplicate mainnet + goerli", input: "mainnet,goerli,mainnet,goerli", wantChains: "mainnet,goerli", wantTargets: "mainnet,goerli"},
		{name: "Preserve target order for non-mainnet", input: "goerli,sepolia", wantChains: "mainnet,goerli,sepolia", wantTargets: "goerli,sepolia"},
		{name: "Case sensitivity preserved", input: "Goerli,goerli", wantChains: "mainnet,Goerli,goerli", wantTargets: "Goerli,goerli"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chains, targets, err := CleanChainString(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantChains, chains)
			assert.Equal(t, tt.wantTargets, targets)
		})
	}
}

func TestCleanChainString_Errors(t *testing.T) {
	t.Run("Internal whitespace error", func(t *testing.T) {
		_, _, err := CleanChainString("goe rli")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInternalWhitespace.Error())
	})
	t.Run("Invalid character error", func(t *testing.T) {
		_, _, err := CleanChainString("goerli!*")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidCharacter.Error())
	})
	t.Run("Empty result error", func(t *testing.T) {
		_, _, err := CleanChainString(",,,")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrEmptyResult.Error())
	})
}
