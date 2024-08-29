package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lucitez/game-alerts/internal/models"
)

type Database struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) Database {
	return Database{conn: conn}
}

func (d Database) GetSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	query := `
		SELECT c.name, c.email, s.id, s.league_id, s.season_id, s.team_name FROM coaches c
		JOIN subscriptions s ON c.id=s.coach_id
		WHERE s.active=true;
	`

	rows, err := d.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := []models.Subscription{}

	for rows.Next() {
		var name string
		var email string
		var id int
		var leagueID string
		var seasonID string
		var teamName string

		if err := rows.Scan(&name, &email, &id, &leagueID, &seasonID, &teamName); err != nil {
			slog.Error("error scanning row", "error", err)
			continue
		}

		subscription := models.Subscription{
			Coach: models.Coach{
				Name:  name,
				Email: email,
			},
			ID:       id,
			LeagueID: leagueID,
			SeasonID: seasonID,
			TeamName: teamName,
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (d Database) CreateSentAlert(ctx context.Context, subscriptionID int, date time.Time) error {
	query := `
		INSERT INTO sent_alerts (subscription_id, game_date)
		VALUES ($1, $2);
	`

	_, err := d.conn.Query(ctx, query, subscriptionID, date)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) HasSentAlert(ctx context.Context, subscriptionID int, date time.Time) (bool, error) {
	query := `
		SELECT game_date FROM sent_alerts
		WHERE subscription_id = $1
		AND game_date = $2;
	`

	row := d.conn.QueryRow(ctx, query, subscriptionID, date)
	var gameDate time.Time
	err := row.Scan(&gameDate)
	if err == nil {
		return true, nil
	}
	if err == pgx.ErrNoRows {
		return false, nil
	}

	return false, err
}
