package wizard

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

// HelpHandlerFunc is a function that provides context-sensitive help
// based on the current screen and question
type HelpHandlerFunc func(screen *Screen, question *Question) string

// Global help handler
var globalHelpHandler HelpHandlerFunc

// SetHelpHandler sets the global help handler function
func SetHelpHandler(handler HelpHandlerFunc) {
	globalHelpHandler = handler
}

// GetHelp returns help text for the current context
func GetHelp(screen *Screen, question *Question) string {
	if globalHelpHandler != nil {
		return globalHelpHandler(screen, question)
	}
	return "Help is not available for this item."
}

type Wizard struct {
	screens    []Screen
	title      string
	current    int
	completed  bool
	displayFn  func(*Wizard, int) error
	Backing    any
	ReloaderFn func(string) (any, error)
}

func NewWizard(screens []Screen, title string, data any, reloaderFn func(string) (any, error)) *Wizard {
	if len(screens) == 0 {
		panic("screens cannot be empty")
	}
	if data == nil || reloaderFn == nil {
		panic("neither data nor reloadFn may be nil")
	}
	ret := &Wizard{
		screens:    screens,
		title:      title,
		current:    0,
		Backing:    data,
		ReloaderFn: reloaderFn,
		displayFn:  displayScreen, // Set the default display function
	}

	// Get screen titles for navigation
	screenTitles := make([]string, len(screens))
	for i, screen := range screens {
		screenTitles[i] = screen.Title
	}

	// Update wizard with navigation bars
	for i := range ret.screens {
		// Create a navigation bar for this screen
		navBar := NewNavigationBar(i, len(screens), screenTitles)

		// Set this wizard as the screen's wizard (for context)
		ret.screens[i].Wizard = ret

		// Set the navigation bar for this screen
		ret.screens[i].NavigationBar = navBar
	}

	return ret
}

func (w *Wizard) Reload(fn string) (err error) {
	w.Backing, err = w.ReloaderFn(fn)
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
		width = max(width, len(screen.Title)+4)
	}
	format := fmt.Sprintf("%s%%-%d.%ds%s\n", colors.Green, width, width, colors.Off)
	for _, screen := range w.screens {
		fmt.Printf(format, screen.Title)
		for _, question := range screen.Questions {
			text, resp := question.GetQuestion()
			if len(text) > 0 {
				cleanText := strings.ReplaceAll(text, "|", "")
				fmt.Printf("  - %s: %s\n", cleanText, resp)
			}
		}
	}

	return nil
}
