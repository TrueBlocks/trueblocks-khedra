package wizard

import (
	"fmt"
	"os"
	"strings"

	coreUtils "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/boxes"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

type Screen struct {
	Title        string
	Subtitle     string
	Body         string
	Instructions string
	Replacements []Replacement
	Questions    []Question
	Style        Style
	Current      int
	Wizard       *Wizard
}

func AddScreen(screen Screen) Screen {
	screen.Title = strings.Trim(screen.Title, "\n")
	screen.Subtitle = strings.Trim(screen.Subtitle, "\n")
	screen.Body = strings.Trim(screen.Body, "\n")
	screen.Instructions = strings.Trim(screen.Instructions, "\n")
	if len(screen.Questions) == 0 {
		screen.Questions = []Question{{}}
	}

	for _, rep := range screen.Replacements {
		screen.Title = rep.Replace(screen.Title)
		screen.Subtitle = rep.Replace(screen.Subtitle)
		screen.Body = rep.Replace(screen.Body)
		screen.Instructions = rep.Replace(screen.Instructions)
		for i := range screen.Questions {
			question := &screen.Questions[i]
			question.Text = strings.ReplaceAll(question.Text, "\n\t\t", "\n          ")
			question.Hint = strings.ReplaceAll(question.Hint, "\n\t\t", "\n          ")
			question.Text = rep.Replace(question.Text)
			for _, rrep := range question.Replacements {
				question.Text = rrep.Replace(question.Text)
			}
		}
	}
	return screen
}

func (s *Screen) OpenHelp() {
	helpMap := map[string]string{
		"KHEDRA WIZARD":     "welcome",
		"General Settings":  "general",
		"Services Settings": "services",
		"Chains Settings":   "chains",
		"Summary":           "summary",
	}
	title := utils.StripColors(s.Title)
	url := "https://khedra.trueblocks.io/user_manual/wizard/" + helpMap[title] + ".html"
	coreUtils.System("open " + url)
}

func (s *Screen) Display(question *Question, caret string) {
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
			b := boxes.Box(lines, screenWidth, boxes.Single|boxes.BottomBorder|boxes.LeftBorder|boxes.RightBorder|boxes.TCorners, style.Justify)
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

	lines := []string{}
	lines = append(lines, titleRows(s.Title, s.Subtitle, &s.Style)...)
	if len(question.Text) > 0 {
		lines = append(lines, question.getLines()...)
	} else {
		lines = append(lines, s.Body)
	}
	lines = append(lines, heightPad(strings.Join(lines, "\n"), 13)...)
	lines = append(lines, s.Instructions)
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
