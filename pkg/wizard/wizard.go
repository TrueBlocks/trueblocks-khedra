package wizard

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

type Wizard struct {
	screens   []Screen
	caret     string
	current   int
	completed bool
	displayFn func(*Wizard, int) error
	Backing   any
	ReloadFn  func(string) (any, error)
}

func NewWizard(screens []Screen, caret string, backing any, reloadFn func(string) (any, error)) *Wizard {
	if len(screens) == 0 {
		panic("screens cannot be empty")
	}
	if backing == nil || reloadFn == nil {
		panic("neither backing nor reloadFn may be nil")
	}

	if caret == "" {
		caret = "--> "
	}

	return &Wizard{
		screens:   screens,
		caret:     caret,
		displayFn: displayScreen,
		Backing:   backing,
		ReloadFn:  reloadFn,
	}
}

func (w *Wizard) Reload(fn string) (err error) {
	w.Backing, err = w.ReloadFn(fn)
	return
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

	fmt.Printf("%s\n", clearScreen)
	fmt.Println("Your answers:")
	width := 0
	for _, screen := range w.screens {
		width = base.Max(width, len(screen.Title)+4)
	}
	format := fmt.Sprintf("%s%%-%d.%ds%s\n", colors.Green, width, width, colors.Off)
	for _, screen := range w.screens {
		fmt.Printf(format, screen.Title)
		for _, question := range screen.Questions {
			text, resp := question.GetQuestion()
			if len(text) > 0 {
				fmt.Printf("  - %s: %s\n", text, resp)
			}
		}
	}

	return nil
}
