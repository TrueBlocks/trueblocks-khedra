package wizard

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type Questioner interface {
	HandleResponse(string) error
	Prompt(string, string, ...bool) string
	GetLines() []string
	Prepare(*Screen) bool
	Clean(*Replacement)
	Clear()
	GetQuestion() (string, string)
	GetError() string
}

// Question models an interactive user prompt.
// Fields:
// - Question: The question displayed to the user.
// - Value: A processed or validated version of the response.
// - Response: An error message displayed in case of invalid input.
// - ErrorStr: An error message displayed in case of invalid input.
// - Prepare: A function for pre-question processing.
// - Validate: A function to validate user input, returning the processed value or an error.
type Question struct {
	Question       string
	Hint           string
	Value          string
	State          string
	Response       string
	ErrorStr       string
	PrepareFn      func(string, *Question) (string, error)
	Validate       func(string, *Question) (string, error)
	Replacements   []Replacement
	Messages       []string
	Screen         *Screen
	ValidationType string // Type of validation to apply ("rpc", "folder", etc.)
}

func (q *Question) HandleResponse(input string) error {
	q.Response = ""
	q.ErrorStr = ""
	input = strings.TrimSpace(input)
	if input == "" {
		input = utils.StripColors(q.Value)
	}

	if q.ValidationType != "" && input != "" {
		if validationFunc := GetValidationFunc(q.ValidationType); validationFunc != nil {
			feedback := validationFunc(input)
			if !feedback.IsValid {
				q.ErrorStr = feedback.Message
				return ErrValidate
			}
			q.Response = FormatValidationFeedback(feedback)
		}
	}

	switch input {
	case "h", "help":
		return ErrUserHelp
	case "q", "quit":
		return ErrUserQuit
	case "e", "edit":
		return ErrUserEdit
	case "c", "chains":
		return ErrUserChains
	case "b", "back":
		return ErrUserBack
	default:
		if q.Validate != nil {
			var err error
			if q.Value, err = q.Validate(input, q); err != nil {
				q.ErrorStr = ""
				q.Response = ""
				if errors.Is(err, ErrValidateWarn) || errors.Is(err, ErrValidateMsg) {
					q.Response = err.Error()
				} else {
					q.ErrorStr = err.Error()
				}
			}
			return err
		} else {
			q.Value = input
			q.Response = input
		}
	}
	return nil
}

func (q *Question) Prompt(str, spacer string, pad ...bool) string {
	if len(pad) > 0 && !pad[0] {
		str = spacer + str
	} else {
		str = spacer + fmt.Sprintf("%-*s", 10, str+":")
	}

	var reps = Replacement{Color: colors.Green, Values: []string{
		"Question:", "Current:", "Answer:", "Error:", "Hint", "Response", "State",
	}}
	return reps.Replace(str)
}

func (q *Question) GetLines() []string {
	var lines []string
	q.Clean(nil)
	if q.Question != "" {
		if len(q.Question) > 0 {
			lines = append(lines, splitLines(q.Prompt("Question", ""), q.Question)...)
		}
		if len(q.Hint) > 0 {
			lines = append(lines, splitLines(q.Prompt("Hint", ""), q.Hint)...)
		}

		if q.State != "" {
			lines = append(lines, q.Prompt("State", "")+colors.Yellow+q.State+colors.Off)
		}
		if q.Value != "" {
			lines = append(lines, q.Prompt("Current", "")+colors.BrightBlue+q.Value+colors.Off)
		}
		if len(q.ErrorStr) > 0 {
			msg := colors.Red + q.ErrorStr + colors.Off
			lines = append(lines, q.Prompt("Error", "")+msg)
			q.ErrorStr = ""
		}
		if len(q.Response) > 0 {
			msg := colors.BrightBlue + q.Response + colors.Off
			lines = append(lines, q.Prompt("Response", "")+msg)
			q.Response = ""
		}
	}
	return append(lines, "")
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

func (q *Question) Clear() {
	q.Value = ""
	q.Response = ""
	q.ErrorStr = ""
}

func (q *Question) Clean(rep *Replacement) {
	if rep != nil {
		q.Question = rep.Replace(q.Question)
		q.Hint = rep.Replace(q.Hint)
	}
	for _, qRep := range q.Replacements {
		q.Question = qRep.Replace(q.Question)
		q.Hint = qRep.Replace(q.Hint)
	}
}

func (q *Question) GetQuestion() (string, string) {
	t := strings.ReplaceAll(q.Question, "\n", " ")
	t = strings.ReplaceAll(t, "           ", " ")
	t = strings.TrimSpace(t)
	r := colors.Magenta + q.Value + colors.Off
	return t, r
}

func (q *Question) GetError() string {
	return q.ErrorStr
}

func splitLines(prompt, line string) []string {
	parts := strings.Split(line, "|")
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if i == 0 {
			parts[i] = prompt + part
		} else {
			parts[i] = strings.Repeat(" ", 10) + part
		}
	}
	return parts
}
