package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/mattn/go-runewidth"
)

type BorderStyle int

const (
	NoBorder BorderStyle = iota
	Single
	Double
)

type BorderPos int

const (
	TopLeft BorderPos = iota
	TopRight
	BottomLeft
	BottomRight
	Horizontal
	Vertical
)

var tokens = map[BorderStyle]map[BorderPos]rune{
	Single: {
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
		Horizontal:  '─',
		Vertical:    '│',
	},
	Double: {
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
		Horizontal:  '═',
		Vertical:    '║',
	},
}

var ErrUserQuit = errors.New("user quit")

func (w *Wizard) showScreen(screenWidth int) (string, error) {
	topBorder := func(width int, bs BorderStyle) []string {
		return []string{
			string(tokens[bs][TopLeft]) + strings.Repeat(string(tokens[bs][Horizontal]), width-2) + string(tokens[bs][TopRight]),
		}
	}

	bottomBorder := func(width int, bs BorderStyle) []string {
		return []string{
			string(tokens[bs][BottomLeft]) + strings.Repeat(string(tokens[bs][Horizontal]), width-2) + string(tokens[bs][BottomRight]),
		}
	}

	pad := func(line string, width int, bs BorderStyle, just Justification) string {
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
			padRight = 0
			// log.Fatalf("line too long in showScreen: %s%s%s [%d,%d,%d,%d]", colors.Red, line, colors.Off, width, lineLen, padLeft, padRight)
		}
		return strings.Repeat(" ", padLeft) + line + strings.Repeat(" ", padRight)
	}

	row := func(s string, width int, bs BorderStyle, just Justification) string {
		body := []string{}
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			padded := pad(line, width, bs, just)
			l := string(tokens[bs][Vertical]) + padded + string(tokens[bs][Vertical])
			body = append(body, l)
		}
		return strings.Join(body, "\n")
	}

	box := func(strs []string, width int, bs BorderStyle, just Justification) string {
		ret := []string{}
		if bs != NoBorder {
			ret = append(ret, topBorder(width, bs)...)
		}
		for _, s := range strs {
			ret = append(ret, row(s, width, bs, just))
		}
		if bs != NoBorder {
			ret = append(ret, bottomBorder(width, bs)...)
		}
		screen := strings.Join(ret, "\n")
		return screen
	}

	titleRows := func(t, s string) []string {
		var ret []string
		if w.Current().Opts.TitleJustify == Center {
			lines := []string{t, "", s}
			ret = []string{box(lines, 34, w.innerBorder, w.Current().Opts.TitleJustify)}
		} else {
			lines := []string{t + " (" + s + ")"}
			b := box(lines, screenWidth, NoBorder, w.Current().Opts.TitleJustify)
			b = strings.TrimSpace(b)
			ret = []string{b}
		}
		ret = append(ret, "")
		return ret
	}

	lines := []string{}
	lines = append(lines, titleRows(w.Current().Title, w.Current().Subtitle)...)
	lines = append(lines, w.Current().Body)
	clearScreen := "\033[2J\033[H"

	fmt.Printf("%s%s\n%s", clearScreen, box(lines, screenWidth, w.outerBorder, Left), w.caret)

	reader := bufio.NewReader(os.Stdin)
	if input, err := reader.ReadString('\n'); err != nil {
		return "", err
	} else {
		input = strings.TrimSpace(input)
		if input == "quit" || input == "q" {
			return "", ErrUserQuit
		}
		if input == "" {
			input = w.Current().Opts.Default
		}
		return input, nil
	}
}
