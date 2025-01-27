package wizard

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

type Replacement struct {
	Color  string
	Values []string
}

func (r *Replacement) Replace(in string) string {
	out := in
	for _, repStr := range r.Values {
		out = strings.ReplaceAll(out, repStr, r.Color+repStr+colors.Off)
	}
	return out
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
