package wizard

type Style struct {
	Outer   Border
	Inner   Border
	Justify Justification
}

func NewStyle() Style {
	return Style{
		Outer:   Single,
		Inner:   Double,
		Justify: Left,
	}
}

type Border int

const (
	NoBorder Border = iota
	Single
	Double
)

type Justification int

const (
	Left Justification = iota
	Right
	Center
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

var boxTokens = map[Border]map[BorderPos]rune{
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
