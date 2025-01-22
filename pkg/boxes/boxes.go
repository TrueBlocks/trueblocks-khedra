package boxes

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
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
		padTotal = base.Max(0, padTotal-1)
	}
	if bs&RightBorder != 0 {
		padTotal = base.Max(0, padTotal-1)
	}
	padLeft, padRight := 0, 0

	switch just {
	case Left:
		padLeft = margin
		padRight = base.Max(0, padTotal-padLeft)
	case Center:
		padLeft = padTotal / 2
		if padTotal%2 != 0 {
			padLeft++
		}
		padRight = base.Max(0, padTotal-padLeft)
	case Right:
		padRight = margin
		padLeft = base.Max(0, padTotal-padRight)
	}

	if padLeft+textWidth+padRight > width {
		padRight = base.Max(0, width-padLeft-textWidth)
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

func Box(strs []string, width int, bs Border, just Justification) string {
	ret := []string{}

	if bs&TopBorder != 0 {
		tb, _ := topBorder(width, bs)
		ret = append(ret, tb...)
	}

	containsAnyRune := func(s string, runes []rune) bool {
		for _, r := range s {
			for _, target := range runes {
				if r == target {
					return true
				}
			}
		}
		return false
	}
	tRunes := []rune{'┬', '├', '┴', '┤', '┼', '╦', '╠', '╩', '╣', '╬'}

	for _, s := range strs {
		if !containsAnyRune(s, tRunes) && bs&(LeftBorder|RightBorder) != 0 {
			ret = append(ret, boxRow(s, width, bs, just))
		} else {
			ret = append(ret, padRow(s, width, bs, just))
		}
	}

	if bs&BottomBorder != 0 {
		bb, _ := bottomBorder(width, bs)
		ret = append(ret, bb...)
	}

	return strings.Join(ret, "\n")
}
