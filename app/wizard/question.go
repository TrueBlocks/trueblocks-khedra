package wizard

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

// Question models an interactive user prompt.
// Fields:
// - Text: The question displayed to the user.
// - Value: A processed or validated version of the response.
// - ErrorMsg: An error message displayed in case of invalid input.
// - Prepare: A function for pre-question processing.
// - Validate: A function to validate user input, returning the processed value or an error.
type Question struct {
	Text         string
	Hint         string
	Value        string
	ErrorMsg     string
	PrepareFn    func(string, *Question) (string, error)
	Validate     func(string, *Question) (string, error)
	Replacements []Replacement
	Messages     []string
	Screen       *Screen
}

func (q *Question) processResponse(input string) error {
	q.ErrorMsg = ""
	input = strings.TrimSpace(input)
	if input == "" {
		input = utils.StripColors(q.Value)
	}
	switch input {
	case "h", "help":
		return ErrUserHelp
	case "q", "quit":
		return ErrUserQuit
	case "b", "back":
		return ErrUserBack
	default:
		if q.Validate != nil {
			var err error
			if q.Value, err = q.Validate(input, q); err != nil {
				q.ErrorMsg = err.Error()
			}
			return err
		} else {
			q.Value = input
		}
	}
	return nil
}

func (q *Question) Prompt(str string, pad ...bool) string {
	var spacer = "  "
	if len(pad) > 0 && !pad[0] {
		str = spacer + str
	} else {
		str = spacer + fmt.Sprintf("%-*s", 10, str+":")
	}

	var reps = Replacement{Color: colors.Green, Values: []string{"Question:", "Current:", "Answer:", "Error:", "Hint"}}
	return reps.Replace(str)
}

func (q *Question) Display(caret string) {
	var lines []string
	if q.Text != "" {
		lines = append(lines, q.Prompt("Question")+q.Text)
		if q.Value != "" {
			value := q.Value
			if len(q.Replacements) > 0 {
				for _, rep := range q.Replacements {
					value = rep.Replace(value)
				}
			}
			lines = append(lines, q.Prompt("Current")+value)
		}
		if q.Hint != "" {
			lines = append(lines, q.Prompt("Hint")+q.Hint)
		}
	}

	if len(q.ErrorMsg) > 0 {
		msg := colors.Red + q.ErrorMsg + colors.Off
		lines = append(lines, q.Prompt("Error")+msg)
		q.ErrorMsg = ""
	}

	lines = append(lines, "")
	fmt.Printf("%s%s", strings.Join(lines, "\n"), q.Prompt(caret, false))
}

func (q *Question) Prepare(s *Screen) bool {
	q.Screen = s
	if q.PrepareFn != nil {
		if value, err := q.PrepareFn(q.Value, q); errors.Is(err, ErrSkipQuestion) {
			return true
		} else {
			q.Value = value
		}
	}

	return false
}
