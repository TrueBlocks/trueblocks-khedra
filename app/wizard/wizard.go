package wizard

import (
	"fmt"
	"log"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

type Justification int

const (
	Left Justification = iota
	Right
	Center
)

type StepType int

const (
	None StepType = iota
	Question
	YesNo
	Screen
)

type Replacement struct {
	Color  string
	Values []string
}

type Option struct {
	Type         StepType
	Default      string
	TitleJustify Justification
	BodyJustify  Justification
	Replacements []Replacement
}

type Step struct {
	Title    string
	Subtitle string
	Body     string
	Response string
	Opts     Option
}

func (s *Step) applyOpts(opts ...Option) {
	if len(opts) == 0 {
		return
	}

	if opts[0].Default != "" {
		s.Opts.Default = opts[0].Default
	}
	if opts[0].Type != None {
		s.Opts.Type = opts[0].Type
	}
	if opts[0].TitleJustify != Left {
		s.Opts.TitleJustify = opts[0].TitleJustify
	}
	if opts[0].BodyJustify != Left {
		s.Opts.BodyJustify = opts[0].BodyJustify
	}
	if len(opts[0].Replacements) > 0 {
		s.Opts.Replacements = opts[0].Replacements
		for _, rep := range opts[0].Replacements {
			for _, val := range rep.Values {
				s.Title = strings.ReplaceAll(s.Title, val, rep.Color+val+colors.Off)
				s.Subtitle = strings.ReplaceAll(s.Subtitle, val, rep.Color+val+colors.Off)
				s.Body = strings.ReplaceAll(s.Body, val, rep.Color+val+colors.Off)
			}
		}
	}
}

func NewStep(title, subtitle, body string, opts ...Option) Step {
	step := Step{
		Title:    strings.Trim(title, "\n"),
		Subtitle: strings.Trim(subtitle, "\n"),
		Body:     strings.Trim(body, "\n"),
		Opts:     Option{Type: Question},
	}
	step.applyOpts(opts...)
	return step
}

func NewScreen(title, subtitle, body string, opts ...Option) Step {
	step := Step{
		Title:    strings.Trim(title, "\n"),
		Subtitle: strings.Trim(subtitle, "\n"),
		Body:     strings.Trim(body, "\n"),
		Opts:     Option{Type: Screen},
	}
	step.applyOpts(opts...)
	return step
}

type Wizard struct {
	step        int
	steps       []Step
	caret       string
	completed   bool
	outerBorder BorderStyle
	innerBorder BorderStyle
}

func NewWizard(steps []Step, caret string) *Wizard {
	if len(steps) == 0 {
		panic("steps cannot be empty")
	}
	if caret == "" {
		caret = "--> "
	}
	return &Wizard{
		steps:       steps,
		caret:       caret,
		outerBorder: Single,
		innerBorder: Double,
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
	// reader := bufio.NewReader(os.Stdin)

	for !w.IsComplete() {
		current := w.Current()
		// body := color(strings.ReplaceAll("{B}"+current.Body, "+", "."))
		// title := color("{W}" + current.Title + "{B}")

		switch current.Opts.Type {
		case Question:
			fallthrough
		case Screen:
			if input, err := w.showScreen(80); err == ErrUserQuit {
				return nil
			} else {
				w.Next(input)
			}
		// case YesNo:
		// 	for {
		// 		fmt.Printf("%s (y/n) %s ", body, w.caret+colors.Off)
		// 		input, _ := reader.ReadString('\n')
		// 		input = strings.TrimSpace(strings.ToLower(input))
		// 		if input == "y" || input == "yes" || input == "n" || input == "no" {
		// 			w.Next(input)
		// 			break
		// 		}
		// 		fmt.Println(colors.Red + "Invalid input. Please enter 'y' or 'n'." + colors.Off)
		// 	}
		default:
			log.Fatalf("unknown step type: %d", current.Opts.Type)
		}
	}

	fmt.Println("\nYour answers:")
	width := 0
	for _, step := range w.GetResponses() {
		width = base.Max(width, len(step.Title)+4)
	}
	format := fmt.Sprintf("%s%%-%d.%ds%s%%s\n", colors.Green, width, width, colors.Off)
	for _, step := range w.GetResponses() {
		fmt.Printf(format, step.Title+":", step.Response)
	}

	return nil
}

func color(in string) string {
	ret := strings.ReplaceAll(in, "{Y}", colors.Yellow)
	ret = strings.ReplaceAll(ret, "{B}", colors.Blue)
	ret = strings.ReplaceAll(ret, "{W}", colors.White)
	return ret
}
