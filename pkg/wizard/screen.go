package wizard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/boxes"
)

type Screen struct {
	Title        string
	Subtitle     string
	Body         string
	Instructions string
	Replacements []Replacement
	Questions    []Questioner
	Style        Style
	Current      int
	Wizard       *Wizard
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
		"KHEDRA WIZARD":     "welcome",
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
	titleRows := func(t, s string, style *Style) []string {
		var ret []string
		if style.Justify == boxes.Center {
			lines := []string{t, "", s}
			ret = []string{boxes.Box(lines, 57, style.Inner, style.Justify)}
		} else {
			lines := []string{t + " (" + s + ")"}
			if len(s) == 0 {
				lines = []string{t}
			}
			lines = append(lines, "──────────────────────────────────────────────────────────────────────────")
			b := boxes.Box(lines, screenWidth-2, boxes.NoBorder, style.Justify)
			b = strings.TrimSpace(b)
			ret = []string{b}
		}
		ret = append(ret, "")
		return ret
	}

	heightPad := func(body string, want int) []string {
		have := len(strings.Split(body, "\n"))
		if have < want {
			return []string{strings.Repeat("\n", want-have-1)}
		}
		return []string{}
	}

	text, _ := question.GetQuestion()

	lines := []string{}
	lines = append(lines, titleRows(s.Title, s.Subtitle, &s.Style)...)
	if len(text) > 0 {
		lines = append(lines, question.GetLines()...)
	} else {
		lines = append(lines, s.Body)
	}
	lines = append(lines, heightPad(strings.Join(lines, "\n"), 13)...)
	if len(text) > 0 {
		lines = append(lines, s.Instructions)
	} else {
		lines = append(lines, "Press enter to continue.")
	}
	screen := boxes.Box(lines, screenWidth, s.Style.Outer, boxes.Left)
	fmt.Printf("%s%s\n%s", clearScreen, screen, question.Prompt(caret, "  ", false))
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
