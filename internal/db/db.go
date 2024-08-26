package db

import (
	"context"
	"log/slog"

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
		SELECT c.name, c.email, s.league_id, s.season_id, s.team_name FROM coaches c
		JOIN subscriptions s ON c.id=s.coach_id
		WHERE s.active=true;
	`

	rows, err := d.conn.Query(ctx, query)
	if err != nil {
		slog.Error("error getting active subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := []models.Subscription{}

	for rows.Next() {
		var name string
		var email string
		var leagueID string
		var seasonID string
		var teamName string

		if err := rows.Scan(&name, &email, &leagueID, &seasonID, &teamName); err != nil {
			slog.Error("error scanning rows", "error", err)
		}

		subscription := models.Subscription{
			Coach: models.Coach{
				Name:  name,
				Email: email,
			},
			LeagueID: leagueID,
			SeasonID: seasonID,
			TeamName: teamName,
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}
