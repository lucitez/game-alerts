package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lucitez/game-alerts/internal/alerter"
	"github.com/lucitez/game-alerts/internal/database"
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
}

func sendGameAlerts(ctx context.Context) error {
	slog.Info("starting send game alerts function")

	slog.Info("getting active subscriptions")
	subscriptions, err := database.GetSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get active subscriptions: %w", err)
	}

	emailer := emailer.New()
	alerter := alerter.New(emailer)

	var alerterErrors error

	slog.Info("sending game alerts")
	for _, subscription := range subscriptions {
		sent, err := alerter.SendGameAlert(ctx, subscription)
		if err != nil {
			alerterErrors = errors.Join(alerterErrors, fmt.Errorf("failed to send game alert: %w", err))
			continue
		}
		if !sent {
			slog.Info("skipped sending game alert", "coach_id", subscription.Coach.ID)
			continue
		}

		slog.Info("sent game alert", "coach_id", subscription.Coach.ID)
	}

	if alerterErrors != nil {
		return alerterErrors
	}

	slog.Info("finished sending game alerts")
	return nil
}
