package wizard

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	coreUtils "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/mattn/go-runewidth"
)

type Screen struct {
	Title        string
	Subtitle     string
	Body         string
	Instructions string
	Replacements []Replacement
	Questions    []Question
	Current      int
	Style        Style
	wiz          *Wizard
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
			question.Text = strings.ReplaceAll(question.Text, "\n\t\t", "\n            ")
			question.Hint = strings.ReplaceAll(question.Hint, "\n\t\t", "\n            ")
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

func (s *Screen) Display() {
	topBorder := func(width int, bs Border) []string {
		return []string{
			string(boxTokens[bs][TopLeft]) + strings.Repeat(string(boxTokens[bs][Horizontal]), width-2) + string(boxTokens[bs][TopRight]),
		}
	}

	bottomBorder := func(width int, bs Border) []string {
		return []string{
			string(boxTokens[bs][BottomLeft]) + strings.Repeat(string(boxTokens[bs][Horizontal]), width-2) + string(boxTokens[bs][BottomRight]),
		}
	}

	padRow := func(line string, width int, bs Border, just Justification) string {
		_ = just // linter
		bw := 2
		padLeft := 3
		if bs == NoBorder {
			bw = 0
			padLeft = 0
			width -= 6
		}
		lineLen := runewidth.StringWidth(utils.StripColors(line))
		padRight := width - bw - lineLen - padLeft
		if padRight < 0 {
			log.Printf("line too long in padRow: %s%s%s [%d,%d,%d,%d]\n", colors.Red, line, colors.Off, width, lineLen, padLeft, padRight)
			padRight = 0
		}
		return strings.Repeat(" ", padLeft) + line + strings.Repeat(" ", padRight)
	}

	boxRow := func(s string, width int, bs Border, just Justification) string {
		body := []string{}
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			padded := padRow(line, width, bs, just)
			l := string(boxTokens[bs][Vertical]) + padded + string(boxTokens[bs][Vertical])
			body = append(body, l)
		}
		return strings.Join(body, "\n")
	}

	box := func(strs []string, width int, bs Border, just Justification) string {
		ret := []string{}

		if bs != NoBorder {
			ret = append(ret, topBorder(width, bs)...)
			for _, s := range strs {
				ret = append(ret, boxRow(s, width, bs, just))
			}
			ret = append(ret, bottomBorder(width, bs)...)
		} else {
			for _, s := range strs {
				ret = append(ret, boxRow(s, width, bs, just))
			}
		}

		return strings.Join(ret, "\n")
	}

	titleRows := func(t, s string, style *Style) []string {
		var ret []string
		if style.Justify == Center {
			lines := []string{t, "", s}
			ret = []string{box(lines, 57, style.Inner, style.Justify)}
		} else {
			lines := []string{t + " (" + s + ")"}
			if len(s) == 0 {
				lines = []string{t}
			}
			b := box(lines, screenWidth, NoBorder, style.Justify)
			b = strings.TrimSpace(b)
			ret = []string{b}
		}
		ret = append(ret, "")
		return ret
	}

	bodyPad := func(body string, want int) []string {
		have := len(strings.Split(body, "\n"))
		if have < want {
			return []string{strings.Repeat("\n", want-have-1)}
		}
		return []string{}
	}
	
	lines := []string{}
	lines = append(lines, titleRows(s.Title, s.Subtitle, &s.Style)...)
	lines = append(lines, s.Body)
	lines = append(lines, bodyPad(strings.Join(lines, "\n"), 13)...)
	lines = append(lines, s.Instructions)
	screen := box(lines, screenWidth, s.Style.Outer, Left)
	clearScreen := "\033[2J\033[H"
	if os.Getenv("NO_CLEAR") == "true" {
		clearScreen = ""
	}
	fmt.Printf("%s%s\n", clearScreen, screen)
}

func (s *Screen) GetCaret(caret string, i, skipped int) string {
	if len(s.Questions) > 1 {
		caret = fmt.Sprintf("%d/%d"+caret, i+1-skipped, len(s.Questions)-skipped)
	}
	return caret

}
