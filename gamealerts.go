package gamealerts

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/lucitez/game-alerts/internal/alerter"
	"github.com/lucitez/game-alerts/internal/db"
	"github.com/lucitez/game-alerts/internal/logger"
)

func init() {
	functions.CloudEvent("SendGameAlerts", sendGameAlerts)

	logger.Init()
}

func sendGameAlerts(ctx context.Context, event cloudevents.Event) error {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to initialize db driver", "error", err)
		return err
	}
	defer conn.Close(context.Background())

	err = conn.Ping(ctx)
	if err != nil {
		slog.Error("failed to start up db", "error", err)
		return err
	}

	db := db.New(conn)

	subscriptions, err := db.GetSubscriptions(ctx)
	if err != nil {
		slog.Error("failed to get active subscriptions", "error", err)
		return err
	}

	a := alerter.New(db)

	for _, subscription := range subscriptions {
		err := a.SendGameAlert(ctx, subscription)
		if err != nil {
			slog.Error("failed to send game alert", "error", err)
		}
		slog.Info("Finished sending game alert", "subscription", subscription)
	}

	slog.Info("Finished sending game alerts")
	return nil
}
