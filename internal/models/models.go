package models

import (
	"database/sql"
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
	SeasonID sql.NullString
	TeamName string
}

func (s Subscription) BuildUrl() string {
	if s.SeasonID.Valid {
		return fmt.Sprintf("https://teampages.com/leagues/%s/events.json?calendar=true&season_id=%s", s.LeagueID, s.SeasonID.String)
	}
	return fmt.Sprintf("https://teampages.com/leagues/%s/events.json?calendar=true", s.LeagueID)
}

func (s Subscription) BuildHumanUrl() string {
	if s.SeasonID.Valid {
		return fmt.Sprintf("https://teampages.com/leagues/%s/events?season_id=%s&view_mode=list", s.LeagueID, s.SeasonID.String)
	}
	return fmt.Sprintf("https://teampages.com/leagues/%s/events?view_mode=list", s.LeagueID)
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
