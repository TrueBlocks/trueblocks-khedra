package wizard

type Border int
type Justification string

const (
	Single Border = 1
	Double Border = 2
	All    Border = 15 // This would be TopBorder | BottomBorder | LeftBorder | RightBorder
)

type Style struct {
	Outer   Border
	Inner   Border
	Justify Justification
}

func NewStyle() Style {
	return Style{
		Outer:   Single | All,
		Inner:   Double | All,
		Justify: "Left",
	}
}
