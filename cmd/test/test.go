package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	ga "github.com/lucitez/game-alerts"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("I am running the thing!")

	fmt.Printf("Env: %s\n", os.Getenv("Env"))

	err := ga.SendGameAlerts(context.Background())
	if err != nil {
		slog.Error("error sending game alerts", "error", err)
	}

	slog.Info("I am done running the thing!")
}
