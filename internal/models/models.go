package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Coach struct {
	Name  string
	Email string
}

type Subscription struct {
	ID       int
	Coach    Coach
	LeagueID string
	SeasonID string
	TeamName string
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

func (g Game) Time() string {
	return fmt.Sprintf(
		"%s, %s %d at %s",
		g.Start.Weekday().String(),
		g.Start.Month().String(),
		g.Start.Day(),
		g.Start.Format(time.Kitchen),
	)
}
