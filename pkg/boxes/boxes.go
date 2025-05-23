package boxes

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/mattn/go-runewidth"
)

type Justification int

const (
	Left Justification = iota
	Right
	Center
)

type Border int

const (
	NoBorder Border = 0
	Single   Border = (1 << iota)
	Double
	TopBorder
	BottomBorder
	LeftBorder
	RightBorder
	TCorners
	Side      = LeftBorder | RightBorder
	TopBottom = TopBorder | BottomBorder
	All       = TopBottom | Side
)

type BorderPos int

const (
	TopLeft BorderPos = iota
	TopRight
	BottomLeft
	BottomRight
	Horizontal
	Vertical
	TopT
	LeftT
	BottomT
	RightT
	MiddleT
)

var boxTokens = map[Border]map[BorderPos]rune{
	Single: {
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
		Horizontal:  '─',
		Vertical:    '│',
		TopT:        '┬',
		LeftT:       '├',
		BottomT:     '┴',
		RightT:      '┤',
		MiddleT:     '┼',
	},
	Double: {
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
		Horizontal:  '═',
		Vertical:    '║',
		TopT:        '╦',
		LeftT:       '╠',
		BottomT:     '╩',
		RightT:      '╣',
		MiddleT:     '╬',
	},
}

func topBorder(width int, bs Border) ([]string, error) {
	key := bs & (Single | Double | NoBorder)
	if width < 2 || boxTokens[key] == nil {
		return nil, fmt.Errorf("invalid width or unsupported border style")
	}
	tokens := boxTokens[key]
	lTok, mTok, rTok := TopLeft, Horizontal, TopRight
	if bs&TCorners != 0 {
		lTok, mTok, rTok = LeftT, Horizontal, RightT
	}
	return []string{string(tokens[lTok]) + strings.Repeat(string(tokens[mTok]), width-2) + string(tokens[rTok])}, nil
}

func bottomBorder(width int, bs Border) ([]string, error) {
	key := bs & (Single | Double | NoBorder)
	if width < 2 || boxTokens[key] == nil {
		return nil, fmt.Errorf("invalid width or unsupported border style")
	}
	tokens := boxTokens[key]
	lTok, mTok, rTok := BottomLeft, Horizontal, BottomRight
	if bs&TCorners != 0 {
		lTok, mTok, rTok = LeftT, Horizontal, RightT
	}
	return []string{string(tokens[lTok]) + strings.Repeat(string(tokens[mTok]), width-2) + string(tokens[rTok])}, nil
}

var margin = 1
var padStr = " "

func padRow(line string, width int, bs Border, just Justification) string {
	textWidth := runewidth.StringWidth(utils.StripColors(line))
	if textWidth >= width {
		return line
	}

	padTotal := width - textWidth
	if bs&LeftBorder != 0 {
		padTotal = max(0, padTotal-1)
	}
	if bs&RightBorder != 0 {
		padTotal = max(0, padTotal-1)
	}
	padLeft, padRight := 0, 0

	switch just {
	case Left:
		padLeft = margin
		padRight = max(0, padTotal-padLeft)
	case Center:
		padLeft = padTotal / 2
		if padTotal%2 != 0 {
			padLeft++
		}
		padRight = max(0, padTotal-padLeft)
	case Right:
		padRight = margin
		padLeft = max(0, padTotal-padRight)
	}

	if padLeft+textWidth+padRight > width {
		padRight = max(0, width-padLeft-textWidth)
	}

	return strings.Repeat(padStr, padLeft) + line + strings.Repeat(padStr, padRight)
}

func boxRow(str string, width int, bs Border, just Justification) string {
	body := []string{}
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		padded := padRow(line, width, bs, just)
		bbs := bs & (Single | Double | NoBorder)
		l := string(boxTokens[bbs][Vertical]) + padded + string(boxTokens[bbs][Vertical])
		body = append(body, l)
	}
	return strings.Join(body, "\n")
}

func CreateBox(strs []string, width int, bs Border, just Justification) string {
	return Box(strs, width, bs, just)
}

func Box(strs []string, width int, bs Border, just Justification) string {
	ret := []string{}

	if bs&TopBorder != 0 {
		tb, _ := topBorder(width, bs)
		ret = append(ret, tb...)
	}

	for _, s := range strs {
		ret = append(ret, boxRow(s, width, bs, just))
	}

	if bs&BottomBorder != 0 {
		bb, _ := bottomBorder(width, bs)
		ret = append(ret, bb...)
	}

	return strings.Join(ret, "\n")
}

// BorderStyle defines the characters used for drawing a box
type BorderStyle struct {
	TopLeft     string
	Top         string
	TopRight    string
	Right       string
	BottomRight string
	Bottom      string
	BottomLeft  string
	Left        string
}

// Style defines the styling options for a box
type Style struct {
	Width       int
	Padding     int
	Justify     Justification
	BorderStyle BorderStyle
}

// NewStyle creates a default box style
func NewStyle() Style {
	return Style{
		Width:   78,
		Padding: 1,
		Justify: Left,
		BorderStyle: BorderStyle{
			TopLeft:     string(boxTokens[Single][TopLeft]),
			Top:         string(boxTokens[Single][Horizontal]),
			TopRight:    string(boxTokens[Single][TopRight]),
			Right:       string(boxTokens[Single][Vertical]),
			BottomRight: string(boxTokens[Single][BottomRight]),
			Bottom:      string(boxTokens[Single][Horizontal]),
			BottomLeft:  string(boxTokens[Single][BottomLeft]),
			Left:        string(boxTokens[Single][Vertical]),
		},
	}
}

// MyBox defines the structure of a box
type MyBox struct {
	Title   string
	Content string
	Style   Style
}

// Display shows the box
func (b *MyBox) Display() {
	// Clear the screen first
	fmt.Print("\033[H\033[2J") // ANSI escape sequence to clear screen and move cursor to home position

	lines := strings.Split(b.Content, "\n")
	width := b.Style.Width

	// If there's a title, add it as the first line
	if b.Title != "" {
		// Add title line and separator
		header := []string{b.Title, strings.Repeat("─", width-4)}
		lines = append(header, lines...)
	}

	// Create and display the box
	boxOutput := Box(lines, width, All, Justification(b.Style.Justify))
	fmt.Println(boxOutput)
}

// NewBox creates a new box with title, content, and style
func NewBox(title, content string, style Style) MyBox {
	return MyBox{
		Title:   title,
		Content: content,
		Style:   style,
	}
}
