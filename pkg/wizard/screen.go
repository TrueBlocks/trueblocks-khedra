package wizard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/boxes"
	"github.com/mattn/go-runewidth"
)

type Screen struct {
	Title         string
	Subtitle      string
	Body          string
	Instructions  string
	Replacements  []Replacement
	Questions     []Questioner
	Style         Style
	Current       int
	Wizard        *Wizard
	NavigationBar *NavigationBar // Add navigation bar field
}

func AddScreen(screen Screen) Screen {
	screen.Title = strings.Trim(screen.Title, "\n")
	screen.Subtitle = strings.Trim(screen.Subtitle, "\n")
	screen.Body = strings.Trim(screen.Body, "\n")
	screen.Instructions = strings.Trim(screen.Instructions, "\n")

	for _, rep := range screen.Replacements {
		screen.Title = rep.Replace(screen.Title)
		screen.Subtitle = rep.Replace(screen.Subtitle)
		screen.Body = rep.Replace(screen.Body)
		screen.Instructions = rep.Replace(screen.Instructions)
		for i := range screen.Questions {
			screen.Questions[i].Clean(&rep)
		}
	}

	return screen
}

func (s *Screen) OpenHelp() {
	helpMap := map[string]string{
		"Welcome Screen":    "welcome",
		"General Settings":  "general",
		"Services Settings": "services",
		"Chain Settings":    "chains",
		"Summary":           "summary",
	}
	title := utils.StripColors(s.Title)
	url := "https://khedra.trueblocks.io/user_manual/wizard/" + helpMap[title] + ".html"
	utils.System("open " + url)
}

func (s *Screen) Display(question Questioner, caret string) {
	heightPad := func(body string, want int) []string {
		have := len(strings.Split(body, "\n"))
		if have < want {
			return []string{strings.Repeat("\n", want-have-1)}
		}
		return []string{}
	}

	questionText, _ := question.GetQuestion()

	// Add navigation bar if available
	var navBarContent string
	if s.NavigationBar != nil {
		navBarContent = s.NavigationBar.Render()
	}

	lines := []string{}

	// Add navigation bar first if available
	if navBarContent != "" {
		lines = append(lines, navBarContent)
	}

	// Show body only when there are no questions or the current question is empty
	if len(questionText) == 0 && len(s.Body) > 0 {
		// Remove leading blank line from body before adding
		body := colors.Green + strings.TrimPrefix(s.Body, "\n") + colors.Off
		lines = append(lines, body)
	}

	// Now add question content if there is any
	if len(questionText) > 0 {
		questionLines := question.GetLines()
		wrappedLines := []string{}

		// Maximum content width to prevent overflow
		boxStyle := boxes.NewStyle()
		maxWidth := boxStyle.Width - 4 // Subtract padding and borders

		// Wrap each line if necessary
		for _, line := range questionLines {
			// Skip empty lines
			if len(line) == 0 {
				wrappedLines = append(wrappedLines, line)
				continue
			}

			// If line is too long, wrap it
			if runewidth.StringWidth(utils.StripColors(line)) > maxWidth {
				// Split by words and reconstruct with wrapping
				words := strings.Fields(line) // Preserve color codes
				var currentLine string

				for _, word := range words {
					// If adding this word would make the line too long, start a new line
					if len(currentLine) > 0 &&
						runewidth.StringWidth(utils.StripColors(currentLine+" "+word)) > maxWidth {
						wrappedLines = append(wrappedLines, currentLine)
						currentLine = word
					} else {
						if len(currentLine) > 0 {
							currentLine += " " + word
						} else {
							currentLine = word
						}
					}
				}

				// Add the last line
				if len(currentLine) > 0 {
					wrappedLines = append(wrappedLines, currentLine)
				}
			} else {
				wrappedLines = append(wrappedLines, line)
			}
		}

		lines = append(lines, wrappedLines...)
	}

	lines = append(lines, heightPad(strings.Join(lines, "\n"), 13)...)
	if len(questionText) > 0 && len(s.Instructions) > 0 {
		lines = append(lines, s.Instructions)
	}

	// Add keyboard shortcuts bar
	shortcutBar := GetShortcutBarForScreen(s.Title, s.Wizard)
	lines = append(lines, shortcutBar)

	// Create a box with updated styles using the new boxes package
	boxStyle := boxes.NewStyle()

	// Create content as string
	content := strings.Join(lines, "\n")

	// Clear the screen first
	fmt.Print("\033[H\033[2J") // ANSI escape sequence to clear screen and move cursor to home position

	// Create a box with single border style
	fmt.Println(boxes.Box(strings.Split(content, "\n"), boxStyle.Width, boxes.Single|boxes.All, boxes.Left))

	// Print the prompt without a newline to avoid extra line
	fmt.Print(question.Prompt(caret+" ", "  ", false))
}

func (s *Screen) GetCaret(caret string, i, skipped int) string {
	if len(s.Questions) > 1 {
		caret = fmt.Sprintf("%d/%d"+caret, i+1-skipped, len(s.Questions)-skipped)
	}
	return caret
}

var clearScreen = "\033[2J\033[H"

func init() {
	if os.Getenv("NO_CLEAR") == "true" {
		clearScreen = ""
	}
}

func (s *Screen) Reload(fn string) error {
	return s.Wizard.Reload(fn)
}

func (s *Screen) EditFile(fn string) error {
	isBlockingEditor := func(editor string) bool {
		if editor == "" {
			return false
		}
		blockingEditors := []string{"nano", "vim", "vi", "emacs -nw", "pico", "ed"}
		for _, be := range blockingEditors {
			if strings.HasPrefix(editor, be) {
				return true
			}
		}
		return false
	}

	editor := os.Getenv("EDITOR")
	if editor == "testing" {
		fmt.Println("Would have edited:")
		return nil
	} else if !isBlockingEditor(editor) {
		editor = "nano"
	}

	args := strings.Split(editor, " ")
	cmd := exec.Command(args[0], append(args[1:], fn)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open config for editing: %w", err)
	}
	return nil
}
