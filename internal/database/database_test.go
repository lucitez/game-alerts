package database

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lucitez/game-alerts/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetSubscriptions(t *testing.T) {
	hs, err := GetSubscriptions()
	assert.NoError(t, err)

	expected := []models.Subscription{
		{
			Coach: models.Coach{
				ID:    1,
				Name:  "Lucas",
				Email: "lucgreggs@gmail.com",
			},
			LeagueID: 150442,
			SeasonID: 2153935,
			TeamName: "Allagash Wanderers",
		},
		{
			Coach: models.Coach{
				ID:    2,
				Name:  "Josh",
				Email: "joshua.lettiere@gmail.com",
			},
			LeagueID: 150442,
			SeasonID: 2153935,
			TeamName: "Allagash Wanderers",
		},
	}

	assert.Equal(t, expected, hs)
}

func TestUpdateSendHistory(t *testing.T) {
	historyPath := "send_history.json"

	originalBytes, err := os.ReadFile(historyPath)
	if !assert.NoError(t, err) {
		return
	}

	t.Cleanup(func() {
		os.WriteFile(historyPath, originalBytes, 0644)
	})

	var historyBefore []historyEntry
	err = json.Unmarshal(originalBytes, &historyBefore)
	if !assert.NoError(t, err) {
		return
	}

	err = UpdateSendHistory(99, 99)
	if !assert.NoError(t, err) {
		return
	}

	afterBytes, err := os.ReadFile(historyPath)
	if !assert.NoError(t, err) {
		return
	}

	var historyAfter []historyEntry
	err = json.Unmarshal(afterBytes, &historyAfter)
	if !assert.NoError(t, err) {
		return
	}

	expected := append(historyBefore, historyEntry{GameID: 99, RecipientID: 99})
	assert.Equal(t, expected, historyAfter)
}
