package wizard

import (
	"errors"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/stretchr/testify/assert"
)

func TestNewScreen(t *testing.T) {
	trimmingTest := func() {
		screen := Screen{
			Title:        "\nHello\n",
			Subtitle:     "\nSubtitle\n",
			Body:         "\nBody\n",
			Instructions: "\nInstructions\n",
		}

		result := AddScreen(screen)

		assert.Equal(t, "Hello", result.Title)
		assert.Equal(t, "Subtitle", result.Subtitle)
		assert.Equal(t, "Body", result.Body)
		assert.Equal(t, "Instructions", result.Instructions)
	}
	t.Run("Trimming Test", func(t *testing.T) { trimmingTest() })

	replacementsTest := func() {
		screen := Screen{
			Title: "Replace this value",
			Replacements: []Replacement{
				{
					Color:  colors.Green,
					Values: []string{"Replace"},
				},
			},
		}

		result := AddScreen(screen)

		expected := colors.Green + "Replace" + colors.Off + " this value"
		assert.Equal(t, expected, result.Title)
	}
	t.Run("Replacements Test", func(t *testing.T) { replacementsTest() })
}

func TestProcessResponse(t *testing.T) {
	validInputTest := func() {
		question := &Question{
			Validate: func(input string, q *Question) (string, error) {
				if input != "valid" {
					return "", errors.New("invalid input")
				}
				return input, nil
			},
		}

		err := question.processResponse("valid")

		assert.NoError(t, err)
		assert.Equal(t, "valid", question.Value)
		assert.Equal(t, "", question.State)
		assert.Equal(t, "", question.ErrorStr)
	}
	t.Run("Valid Input Test", func(t *testing.T) { validInputTest() })

	invalidInputTest := func() {
		question := &Question{
			Validate: func(input string, q *Question) (string, error) {
				if input != "valid" {
					return "", errors.New("invalid input")
				}
				return input, nil
			},
		}

		err := question.processResponse("invalid")

		assert.Error(t, err)
		assert.Equal(t, "invalid input", question.ErrorStr)
		assert.Equal(t, "", question.Value)
		assert.Equal(t, "", question.State)
	}
	t.Run("Invalid Input Test", func(t *testing.T) { invalidInputTest() })

	commandTest := func() {
		question := &Question{}

		err := question.processResponse("help")
		assert.Equal(t, ErrUserHelp, err)

		err = question.processResponse("quit")
		assert.Equal(t, ErrUserQuit, err)

		err = question.processResponse("back")
		assert.Equal(t, ErrUserBack, err)
	}
	t.Run("Command Test", func(t *testing.T) { commandTest() })
}
