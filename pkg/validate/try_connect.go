package validate

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
)

func TryConnect(chain, url string, attempts int) error {
	for i := 1; i <= attempts; i++ {
		err := rpc.PingRpc(url)
		if err == nil {
			return nil
		} else {
			slog.Warn("retrying RPC", "chain", chain, "provider", url)
			if i < attempts {
				time.Sleep(1 * time.Second)
			}
		}
	}

	fv := NewFieldValidator("ping_rpc", "Chain", "rpc", fmt.Sprintf("[%s]", chain))
	return Failed(fv, fmt.Sprintf("cannot connect to RPC (%s-%s) after %d attempts", chain, url, attempts), "")
}
