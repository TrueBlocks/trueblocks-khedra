package app

import (
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

func HasValidRpc(ch *types.Chain, tries int) bool {
	for _, rpc := range ch.RPCs {
		if err := types.TryConnect(ch.Name, rpc, tries); err == nil {
			return true
		}
	}
	return false
}
