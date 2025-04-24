package wizard

import (
	"strings"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestReplacement(t *testing.T) {
	validReplacement := func() {
		defer types.SetupTest([]string{})()
		replacement := &Replacement{
			Color:  colors.Blue,
			Values: []string{"value1", "value2"},
		}

		err := replacement.Validate()
		assert.NoError(t, err)
		assert.Contains(t, replacement.Values, "value1")
	}
	t.Run("Valid Replacement", func(t *testing.T) { validReplacement() })

	missingFields := func() {
		defer types.SetupTest([]string{})()
		replacement := &Replacement{
			Color:  "",
			Values: []string{"value1", "value2"},
		}

		err := replacement.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "color field is empty")
	}
	t.Run("Missing Fields", func(t *testing.T) { missingFields() })

	emptyReplacement := func() {
		defer types.SetupTest([]string{})()
		replacement := &Replacement{
			Color:  "blue",
			Values: []string{},
		}

		err := replacement.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "values field is empty")
	}
	t.Run("Empty Replacement", func(t *testing.T) { emptyReplacement() })

	testReplaceMethod := func() {
		defer types.SetupTest([]string{})()
		replacement := &Replacement{
			Color:  colors.Green,
			Values: []string{"test", "replace"},
		}

		input := "this is a test to replace"
		expected := strings.ReplaceAll(input, "test", colors.Green+"test"+colors.Off)
		expected = strings.ReplaceAll(expected, "replace", colors.Green+"replace"+colors.Off)

		result := replacement.Replace(input)
		assert.Equal(t, expected, result)
	}
	t.Run("Test Replace Method", func(t *testing.T) { testReplaceMethod() })
}
