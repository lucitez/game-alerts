package alerter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/lucitez/game-alerts/internal/db"
	"github.com/lucitez/game-alerts/internal/emailer"
	"github.com/lucitez/game-alerts/internal/models"
)

type Alerter struct {
	db      db.Database
	emailer emailer.Emailer
}

func New(db db.Database, emailer emailer.Emailer) Alerter {
	return Alerter{
		db:      db,
		emailer: emailer,
	}
}

func (a Alerter) SendGameAlert(ctx context.Context, subscription models.Subscription) error {
	nextGame, err := getNextGame(subscription)
	if err != nil {
		return fmt.Errorf("failed to get the next game: %w", err)
	}
	if nextGame == (Game{}) {
		slog.Info("Next game has not been posted yet")
		return nil
	}

	hasSentAlert, err := a.db.HasSentAlert(ctx, subscription.ID, nextGame.Start)
	if err != nil {
		return fmt.Errorf("failed to get hasSentAlert: %w", err)
	}
	if hasSentAlert {
		slog.Info("Already sent game alert", "nextGame", nextGame)
		return nil
	}

	err = a.sendGameAlertEmail(nextGame, subscription)
	if err != nil {
		return fmt.Errorf("failed to send game alert email: %w", err)
	}

	err = a.db.CreateSentAlert(ctx, subscription.ID, nextGame.Start)
	if err != nil {
		return fmt.Errorf("failed to create sent alert: %w", err)
	}

	return nil
}

type Game struct {
	Start    time.Time `json:"start"`
	HomeTeam string    `json:"homeTeam"`
	AwayTeam string    `json:"awayTeam"`
	Location string    `json:"location"`
}

func (g Game) Field() (string, error) {
	re := regexp.MustCompile(`^.+ - (.+)$`)
	matches := re.FindStringSubmatch(g.Location)
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to extract field from game: %+v", g)
	}
	return strings.ToLower(re.FindStringSubmatch(g.Location)[1]), nil
}

func getNextGame(subscription models.Subscription) (Game, error) {
	teamName := subscription.TeamName
	leagueID := subscription.LeagueID
	seasonID := subscription.SeasonID

	url := fmt.Sprintf("https://teampages.com/leagues/%s/events.json?calendar=true&season_id=%s", leagueID, seasonID)
	resp, err := http.Get(url)
	if err != nil {
		return Game{}, err
	}
	defer resp.Body.Close()

	games := []Game{}
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return Game{}, err
	}

	var nextGame Game

	// Games are sorted chronologically. The first game after time.Now() is the next game
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

func (a Alerter) sendGameAlertEmail(game Game, subscription models.Subscription) error {
	field, err := game.Field()
	if err != nil {
		return err
	}

	subject := "Co-ed soccer game"
	body := fmt.Sprintf(
		"Like if you're playing on %s, %s %d at %s (%s field)\r\n"+
			"https://teampages.com/leagues/%s/events?season_id=%s?view_mode=list",
		game.Start.Weekday().String(),
		game.Start.Month().String(),
		game.Start.Day(),
		game.Start.Format(time.Kitchen),
		field,
		subscription.LeagueID,
		subscription.SeasonID,
	)

	err = a.emailer.SendEmail(subscription.Coach.Email, subject, body)
	return err
}
