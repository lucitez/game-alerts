package gamealerts

import (
	"context"
	"log/slog"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/lucitez/game-alerts/internal/alerter"
	"github.com/lucitez/game-alerts/internal/db"
	"github.com/lucitez/game-alerts/internal/emailer"
	"github.com/lucitez/game-alerts/internal/logger"
)

func init() {
	functions.CloudEvent("SendGameAlerts", sendGameAlerts)

	logger.Init()
}

func sendGameAlerts(ctx context.Context, event cloudevents.Event) error {
	slog.Info("starting send game alerts function")

	slog.Info("connecting to db")
	conn, err := db.CreateConnection(ctx)
	if err != nil {
		slog.Error("failed to create database connection", "error", err)
	}

	db := db.New(conn)

	slog.Info("getting active subscriptions")
	subscriptions, err := db.GetSubscriptions(ctx)
	if err != nil {
		slog.Error("failed to get active subscriptions", "error", err)
		return err
	}

	emailer := emailer.New()
	alerter := alerter.New(db, emailer)

	slog.Info("sending game alerts")
	for _, subscription := range subscriptions {
		err := alerter.SendGameAlert(ctx, subscription)
		if err != nil {
			slog.Error("failed to send game alert", "error", err)
			continue
		}
		slog.Info("Finished sending game alert", "subscription", subscription)
	}

	slog.Info("Finished sending game alerts")
	return nil
}
