package alerter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/lucitez/game-alerts/internal/db"
	"github.com/lucitez/game-alerts/internal/emailer"
	"github.com/lucitez/game-alerts/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Emailer interface {
	SendEmail(toEmail string, subject string, body string) error
}

type Alerter struct {
	db      db.Database
	emailer Emailer
}

func New(db db.Database, emailer emailer.Emailer) Alerter {
	return Alerter{
		db:      db,
		emailer: emailer,
	}
}

func (a Alerter) SendGameAlert(ctx context.Context, subscription models.Subscription) (sent bool, err error) {
	nextGame, err := getNextGame(subscription)
	if err != nil {
		return false, fmt.Errorf("failed to get the next game: %w", err)
	}
	if nextGame == (models.Game{}) {
		slog.Info("next game has not been posted yet")
		return false, nil
	}
	if nextGame.Start.After(time.Now().Add(time.Hour * 24 * 8)) {
		slog.Info("next game is more than a week away, holding off for now")
		return false, nil
	}

	hasSentAlert, err := a.db.HasSentAlert(ctx, subscription.ID, nextGame.Start)
	if err != nil {
		return false, fmt.Errorf("failed to get hasSentAlert: %w", err)
	}
	if hasSentAlert {
		slog.Info("already sent game alert", "subscription_id", subscription.ID)
		return false, nil
	}

	err = a.sendEmail(nextGame, subscription)
	if err != nil {
		return false, fmt.Errorf("failed to send game alert email: %w", err)
	}

	err = a.db.CreateSentAlert(ctx, subscription.ID, nextGame.Start)
	if err != nil {
		return false, fmt.Errorf("failed to create sent alert: %w", err)
	}

	return true, nil
}

func getNextGame(subscription models.Subscription) (models.Game, error) {
	teamName := subscription.TeamName
	leagueID := subscription.LeagueID
	seasonID := subscription.SeasonID

	url := fmt.Sprintf("https://teampages.com/leagues/%s/events.json?calendar=true&season_id=%s", leagueID, seasonID)
	resp, err := http.Get(url)
	if err != nil {
		return models.Game{}, err
	}
	defer resp.Body.Close()

	games := []models.Game{}
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return models.Game{}, err
	}

	// Games are supposed to be sorted chronologically, but sometimes they get out of order
	sort.Slice(games, func(i, j int) bool {
		return games[i].Start.Before(games[j].Start)
	})

	var nextGame models.Game
	for _, game := range games {
		if game.HomeTeam != teamName && game.AwayTeam != teamName {
			continue
		}

		if game.Start.After(time.Now()) {
			nextGame = game
			break
		}
	}

	return nextGame, nil
}

func (a Alerter) sendEmail(game models.Game, subscription models.Subscription) error {
	field, err := game.Field()
	if err != nil {
		return err
	}

	c := cases.Title(language.AmericanEnglish)
	sendee := c.String(subscription.Coach.Name)
	if sendee == "" {
		sendee = subscription.Coach.Email
	}

	subject := "Soccer game scheduled"
	body := fmt.Sprintf(
		"Hello %s,\n\nYour next game has been scheduled for %s.\n\nHere's a text blast to copy and paste:\n\n"+
			"Like if you're playing on %s (%s field)\n\n"+
			"TeamPages url: https://teampages.com/leagues/%s/events?season_id=%s?view_mode=list",
		sendee,
		game.Time(),
		game.Time(),
		field,
		subscription.LeagueID,
		subscription.SeasonID,
	)

	err = a.emailer.SendEmail(subscription.Coach.Email, subject, body)
	return err
}
