package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lucitez/game-alerts/internal/alerter"
	"github.com/lucitez/game-alerts/internal/db"
	"github.com/lucitez/game-alerts/internal/emailer"
	"github.com/lucitez/game-alerts/internal/logger"
)

func main() {
	logger.Init()

	err := sendGameAlerts(context.Background())
	if err != nil {
		slog.Error("error sending game alerts", "error", err)
		os.Exit(1)
	}
	os.Exit(1)
}

func sendGameAlerts(ctx context.Context) error {
	slog.Info("starting send game alerts function")

	slog.Info("connecting to db")
	conn, err := db.CreateConnection(ctx)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer conn.Close(ctx)

	db := db.New(conn)

	slog.Info("getting active subscriptions")
	subscriptions, err := db.GetSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active subscriptions: %w", err)
	}

	emailer := emailer.New()
	alerter := alerter.New(db, emailer)

	slog.Info("sending game alerts")
	for _, subscription := range subscriptions {
		sent, err := alerter.SendGameAlert(ctx, subscription)
		if err != nil {
			slog.Error("failed to send game alert", "error", err, "subscription", subscription)
			continue
		}
		if sent {
			slog.Info("sent game alert", "subscription_id", subscription.ID)
			continue
		}
		slog.Info("skipped sending game alert", "subscription_id", subscription.ID)
	}

	slog.Info("finished sending game alerts")
	return nil
}
