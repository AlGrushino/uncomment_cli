package main

import (
	"log/slog"
	"os"
	"uncomment-cli/cmd"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cmd.Execute()
}
