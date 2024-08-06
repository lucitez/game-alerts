package gamealerts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func init() {
	functions.CloudEvent("SendGameAlert", sendGameAlert)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replacer}))
	slog.SetDefault(logger)
}

func replacer(groups []string, a slog.Attr) slog.Attr {
	// Rename attribute keys to match Cloud Logging structured log format
	switch a.Key {
	case slog.LevelKey:
		a.Key = "severity"
		// Map slog.Level string values to Cloud Logging LogSeverity
		// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
		if level := a.Value.Any().(slog.Level); level == slog.LevelWarn {
			a.Value = slog.StringValue("WARNING")
		}
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.MessageKey:
		a.Key = "message"
	}
	return a
}

func sendGameAlert(context.Context, cloudevents.Event) error {
	nextGame, err := getNextGame()
	if err != nil {
		slog.Error("error getting the next game", "error", err)
		return err
	}

	if nextGame == (Game{}) {
		slog.Info("Next game has not been posted yet")
		return errors.New("next game has not been posted yet")
	}

	err = sendEmail(nextGame)
	if err != nil {
		slog.Error("error sending email", "error", err, "game", nextGame)
		return err
	}

	slog.Info("Email send success", "game", nextGame)
	slog.Error("Test error slog", "error", errors.New("test error"))
	return nil
}

type Game struct {
	Start    time.Time `json:"start"`
	HomeTeam string    `json:"homeTeam"`
	AwayTeam string    `json:"awayTeam"`
	Location string    `json:"location"`
}

func (g Game) String() string {
	field, _ := g.Field()

	return fmt.Sprintf("Date: %s; Field: %s", g.Start.String(), field)
}

func (g Game) Field() (string, error) {
	re := regexp.MustCompile(`^.+ - (.+)$`)
	matches := re.FindStringSubmatch(g.Location)
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to extract field from game: %+v", g)
	}
	return strings.ToLower(re.FindStringSubmatch(g.Location)[1]), nil
}

func getNextGame() (Game, error) {
	leagueId := os.Getenv("LEAGUE_ID")
	seasonId := os.Getenv("SEASON_ID")
	url := fmt.Sprintf("https://teampages.com/leagues/%s/events.json?calendar=true&season_id=%s", leagueId, seasonId)
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
	teamName := os.Getenv("TEAM_NAME")

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

// https://pkg.go.dev/net/smtp#example-SendMail
func sendEmail(game Game) error {
	fromEmail := os.Getenv("FROM_EMAIL")
	toEmail := os.Getenv("TO_EMAIL")
	appPass := os.Getenv("APP_PASS")
	leagueId := os.Getenv("LEAGUE_ID")
	seasonId := os.Getenv("SEASON_ID")

	auth := smtp.PlainAuth("", fromEmail, appPass, "smtp.gmail.com")

	field, err := game.Field()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(
		"Like if you're playing on %s, %s %d at %s (%s field)\r\n\r\n"+
			"https://teampages.com/leagues/%s/events?season_id=%s?view_mode=list",
		game.Start.Weekday().String(),
		game.Start.Month().String(),
		game.Start.Day(),
		game.Start.Format(time.Kitchen),
		field,
		leagueId,
		seasonId,
	)

	fullMsg := []byte(
		"Subject: Co-ed soccer game\r\n\r\n" + msg)

	err = smtp.SendMail("smtp.gmail.com:587", auth, fromEmail, []string{toEmail}, fullMsg)
	if err != nil {
		return err
	}

	return nil
}
