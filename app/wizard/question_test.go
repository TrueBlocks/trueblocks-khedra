package wizard

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	initializationTest := func() {
		question := &Question{
			Question:  "What is your name?",
			PrepareFn: func(string, *Question) (string, error) { return "Prepared question", nil },
			Validate: func(input string, q *Question) (string, error) {
				if input == "" {
					return "", errors.New("input is empty")
				}
				return input, nil
			},
		}

		assert.Equal(t, "What is your name?", question.Question)
		assert.Equal(t, "", question.Value)
		assert.Equal(t, "", question.State)
		assert.Equal(t, "", question.ErrorStr) // Not initialized until processResponse
		assert.NotNil(t, question.PrepareFn)
		assert.NotNil(t, question.Validate)
	}
	t.Run("Initialization Test", func(t *testing.T) { initializationTest() })

	processResponseTest := func() {
		question := &Question{
			Question: "Choose an option:",
			Validate: func(input string, q *Question) (string, error) {
				if input != "download" && input != "scratch" {
					return input, fmt.Errorf(`value must be either "download" or "scratch"%w`, ErrValidate)
				}
				return input, nil
			},
		}

		// Valid input
		err := question.processResponse("download")
		assert.NoError(t, err)
		assert.Equal(t, "download", question.Value)
		assert.Equal(t, "", question.State)
		assert.Equal(t, "", question.ErrorStr)

		// Invalid input
		err = question.processResponse("invalid")
		assert.Error(t, err)
		assert.Equal(t, "invalid", question.Value)
		assert.Equal(t, "", question.State)
		assert.Contains(t, question.ErrorStr, "value must be either \"download\" or \"scratch\"")

		// Empty input (uses current Value)
		question.Value = "scratch"
		err = question.processResponse("")
		assert.NoError(t, err)
		assert.Equal(t, "scratch", question.Value)
		assert.Equal(t, "", question.State)
		assert.Equal(t, "", question.ErrorStr)
	}
	t.Run("Process Response Test", func(t *testing.T) { processResponseTest() })

	edgeCasesTest := func() {
		emptyQuestion := &Question{
			Question:  "",
			Value:     "",
			State:     "",
			ErrorStr:  "",
			PrepareFn: nil,
			Validate:  nil,
		}

		assert.Equal(t, "", emptyQuestion.Question)
		assert.Equal(t, "", emptyQuestion.Value)
		assert.Equal(t, "", emptyQuestion.State)
		assert.Equal(t, "", emptyQuestion.ErrorStr)
		assert.Nil(t, emptyQuestion.PrepareFn)
		assert.Nil(t, emptyQuestion.Validate)
	}
	t.Run("Edge Cases Test", func(t *testing.T) { edgeCasesTest() })
}
