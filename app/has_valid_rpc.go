package app

import (
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/validate"
)

func HasValidRpc(ch *types.Chain, tries int) bool {
	for _, rpc := range ch.RPCs {
		if err := validate.TryConnect(ch.Name, rpc, tries); err == nil {
			return true
		}
	}
	return false
}
