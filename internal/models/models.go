package models

type Coach struct {
	Name  string
	Email string
}

type Subscription struct {
	Coach    Coach
	LeagueID string
	SeasonID string
	TeamName string
}
