package gamealerts

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/lucitez/game-alerts/internal/emailer"
	"github.com/lucitez/game-alerts/internal/logger"
)

func init() {
	functions.CloudEvent("SendGameAlert", sendGameAlert)

	logger.Init()
}

func sendGameAlert(context.Context, cloudevents.Event) error {
	nextGame, err := getNextGame()
	if err != nil {
		slog.Error("error getting the next game", "error", err)
		return err
	}
	if nextGame == (Game{}) {
		slog.Info("Next game has not been posted yet")
		return nil
	}

	err = sendGameAlertEmail(nextGame)
	if err != nil {
		slog.Error("error sending email", "error", err, "game", nextGame)
		return err
	}

	slog.Info("Email send success", "game", nextGame)
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

// TODO get these from db
const LEAGUE_ID = "54397"
const SEASON_ID = "2146954"
const TEAM_NAME = "HOW I MEGGED YOUR MOTHER"

func getNextGame() (Game, error) {
	url := fmt.Sprintf("https://teampages.com/leagues/%s/events.json?calendar=true&season_id=%s", LEAGUE_ID, SEASON_ID)
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
		if game.HomeTeam != TEAM_NAME && game.AwayTeam != TEAM_NAME {
			continue
		}

		if game.Start.After(time.Now()) {
			nextGame = game
			break
		}
	}

	return nextGame, nil
}

func sendGameAlertEmail(game Game) error {
	toEmail := os.Getenv("TO_EMAIL")

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
		LEAGUE_ID,
		SEASON_ID,
	)

	err = e.SendEmail(toEmail, subject, body)
	return err
}