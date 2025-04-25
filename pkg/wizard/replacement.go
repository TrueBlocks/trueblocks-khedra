package wizard

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

// Replacement defines a color replacement for text
type Replacement struct {
	Color  string
	Values []string
}

// Replace applies the color replacement to the given text
func (r *Replacement) Replace(text string) string {
	for _, val := range r.Values {
		text = strings.ReplaceAll(text, val, r.Color+val+colors.Off)
	}
	return text
}

func (r *Replacement) Validate() error {
	if r.Color == "" {
		return fmt.Errorf("color field is empty")
	}
	if len(r.Values) == 0 {
		return fmt.Errorf("values field is empty")
	}
	return nil
}
