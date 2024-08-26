package alerter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/lucitez/game-alerts/internal/emailer"
	"github.com/lucitez/game-alerts/internal/models"
)

func SendGameAlert(subscription models.Subscription) error {
	nextGame, err := getNextGame(subscription)
	if err != nil {
		slog.Error("error getting the next game", "error", err)
		return err
	}
	if nextGame == (Game{}) {
		slog.Info("Next game has not been posted yet")
		return nil
	}

	err = sendGameAlertEmail(nextGame, subscription)
	if err != nil {
		slog.Error("error sending email", "error", err, "game", nextGame)
		return err
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

	// Games are sorted chronologically. The first game after now() is the next game
	for _, game := range games {
		if game.HomeTeam != teamName && game.AwayTeam != teamName {
			continue
		}

		if game.Start.After(time.Now().Add(-time.Hour * 24 * 30)) {
			nextGame = game
			break
		}
	}

	return nextGame, nil
}

func sendGameAlertEmail(game Game, subscription models.Subscription) error {
	field, err := game.Field()
	if err != nil {
		return err
	}

	e := emailer.New()

	subject := "Co-ed soccer game"
	body := fmt.Sprintf(
		"Like if you're playing on %s, %s %d at %s (%s field)\r\n\r\n"+
			"https://teampages.com/leagues/%s/events?season_id=%s?view_mode=list",
		game.Start.Weekday().String(),
		game.Start.Month().String(),
		game.Start.Day(),
		game.Start.Format(time.Kitchen),
		field,
		subscription.LeagueID,
		subscription.SeasonID,
	)

	err = e.SendEmail(subscription.Coach.Email, subject, body)
	return err
}
