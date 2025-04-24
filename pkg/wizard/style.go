package wizard

import "github.com/TrueBlocks/trueblocks-khedra/v5/pkg/boxes"

type Style struct {
	Outer   boxes.Border
	Inner   boxes.Border
	Justify boxes.Justification
}

func NewStyle() Style {
	return Style{
		Outer:   boxes.Single | boxes.All,
		Inner:   boxes.Double | boxes.All,
		Justify: boxes.Left,
	}
}
