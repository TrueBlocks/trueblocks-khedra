package app

import (
	"context"
	"log"
	"log/slog"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func Progress(msg string, args ...any) {
	slog.Log(context.Background(), types.LevelProgress, msg, args...)
}

func (k *KhedraApp) Fatal(msg string) {
	log.Fatal(msg)
}
