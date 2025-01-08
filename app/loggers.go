package app

import (
	"fmt"
	"log"
	"os"
	"time"
)

func (k *KhedraApp) Debug(msg string, v ...any) {
	k.fileLogger.Debug(msg, v...)
}

func (k *KhedraApp) Info(msg string, v ...any) {
	k.fileLogger.Info(msg, v...)
	k.progLogger.Info(msg, v...)
}

func (k *KhedraApp) Warn(msg string, v ...any) {
	k.fileLogger.Warn(msg, v...)
	k.progLogger.Warn(msg, v...)
}

func (k *KhedraApp) Error(msg string, v ...any) {
	k.fileLogger.Error(msg, v...)
	k.progLogger.Error(msg, v...)
}

func (k *KhedraApp) Prog(msg string, v ...any) {
	if len(v) > 0 && fmt.Sprint(v[len(v)-1]) == "\n" {
		k.progLogger.Info(msg, v...)
	} else {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		message := fmt.Sprintf("PROG %s %s: %s", timestamp, msg, fmt.Sprint(v...))
		fmt.Fprintf(os.Stdout, "%s\r", message)
		os.Stdout.Sync()
	}
}

func (k *KhedraApp) Fatal(msg string) {
	log.Fatal(msg)
	// k.progLogger.Fatal(msg, v...)
	os.Exit(1)
}
