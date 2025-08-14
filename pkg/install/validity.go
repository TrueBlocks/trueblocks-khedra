package install

import (
	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
)

func Configured() bool {
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return false
	}
	return true
}
