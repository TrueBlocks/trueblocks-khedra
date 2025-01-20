package wizard

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

type Wizard struct {
	screens   []Screen
	caret     string
	current   int
	completed bool
	displayFn func(*Wizard, int) error
}

func NewWizard(screens []Screen, caret string) *Wizard {
	if len(screens) == 0 {
		panic("screens cannot be empty")
	}

	if caret == "" {
		caret = "--> "
	}

	return &Wizard{
		screens:   screens,
		caret:     caret,
		displayFn: displayScreen,
	}
}

func (w *Wizard) Current() *Screen {
	if w.current < 0 || w.current >= len(w.screens) {
		return &Screen{}
	}
	return &w.screens[w.current]
}

func (w *Wizard) Next() bool {
	if w.current+1 >= len(w.screens) {
		w.completed = true
		return false
	}
	w.current++
	w.screens[w.current].Current = 0
	return true
}

func (w *Wizard) Prev() bool {
	if w.current <= 0 {
		return true
	}
	w.current--
	w.screens[w.current].Current = len(w.screens[w.current].Questions) - 1
	return true
}

func (w *Wizard) IsComplete() bool {
	return w.completed
}

func (w *Wizard) Run() error {
	for !w.IsComplete() {
		if err := w.displayFn(w, w.current); err == ErrUserQuit {
			return nil
		} else if err == ErrUserBack {
			if !w.Prev() {
				return nil
			}
		}
	}

	fmt.Println("")
	fmt.Println("Your answers:")
	width := 0
	for _, screen := range w.screens {
		width = base.Max(width, len(screen.Title)+4)
	}
	format := fmt.Sprintf("%s%%-%d.%ds%s\n", colors.Green, width, width, colors.Off)
	for _, screen := range w.screens {
		fmt.Printf(format, screen.Title)
		for _, question := range screen.Questions {
			if len(question.Text) == 0 {
				continue
			}
			fmt.Printf("  - %s: %s\n", strings.TrimSpace(strings.ReplaceAll(question.Text, "            ", "      ")), colors.Magenta+question.Value+colors.Off)
		}
	}

	return nil
}
