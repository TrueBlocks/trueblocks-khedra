package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

type Step struct {
	Name     string
	Type     string
	Prompt   string
	Response string
	Metadata map[string]interface{}
}

type Wizard struct {
	step      int
	steps     []Step
	caret     string
	completed bool
}

func NewWizard(steps []Step, caret string) *Wizard {
	if len(steps) == 0 {
		panic("steps cannot be empty")
	}
	if caret == "" {
		caret = ">"
	}
	return &Wizard{
		steps: steps,
		caret: caret,
	}
}

func (w *Wizard) Step() int {
	return w.step
}

func (w *Wizard) Current() Step {
	if w.step < 0 || w.step >= len(w.steps) {
		return Step{}
	}
	return w.steps[w.step]
}

func (w *Wizard) Next(resp string) bool {
	w.steps[w.step].Response = strings.TrimSpace(resp)
	if w.step+1 >= len(w.steps) {
		w.completed = true
		return false
	}
	w.step++
	return true
}

func (w *Wizard) Prev() bool {
	if w.step <= 0 {
		return false
	}
	w.step--
	return true
}

func (w *Wizard) IsComplete() bool {
	return w.completed
}

func (w *Wizard) Reset() {
	w.step = 0
	w.completed = false
	for i := range w.steps {
		w.steps[i].Response = ""
	}
}

func (w *Wizard) GetResponses() []Step {
	return w.steps
}

func (w *Wizard) Run() error {
	reader := bufio.NewReader(os.Stdin)

	for !w.IsComplete() {
		current := w.Current()

		switch current.Type {
		case "welcome":
			fmt.Println(colors.Blue + current.Prompt + colors.Off)
			fmt.Print("Press Enter to continue...")
			_, _ = reader.ReadString('\n')
			w.Next("")
		case "yes/no":
			for {
				fmt.Printf("%s (y/n) %s ", colors.Blue+current.Prompt, w.caret+colors.Off)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(strings.ToLower(input))

				if input == "y" || input == "yes" || input == "n" || input == "no" {
					w.Next(input)
					break
				}
				fmt.Println(colors.Red + "Invalid input. Please enter 'y' or 'n'." + colors.Off)
			}
		default:
			defaultAnswer := ""
			if val, ok := current.Metadata["default"]; ok {
				defaultAnswer, _ = val.(string)
			}

			fmt.Printf("%s (%s) %s ", colors.Blue+current.Prompt, defaultAnswer, w.caret+colors.Off)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "quit" {
				return errors.New("wizard aborted by user")
			}

			if input == "" {
				input = defaultAnswer
			}

			w.Next(input)
		}
	}

	fmt.Println("\nYour answers:")
	width := 0
	for _, step := range w.GetResponses() {
		width = base.Max(width, len(step.Name)+4)
	}
	format := fmt.Sprintf("%s%%-%d.%ds%s%%s\n", colors.Green, width, width, colors.Off)
	for _, step := range w.GetResponses() {
		fmt.Printf(format, step.Name+":", step.Response)
	}

	return nil
}
