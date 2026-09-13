package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucitez/game-alerts/internal/models"
)

type subscription struct {
	LeagueID   int    `json:"league_id"`
	SeasonID   int    `json:"season_id"`
	TeamName   string `json:"team_name"`
	Recipients []int  `json:"recipients"`
}

type recipient struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func getAbsPath(relPath string) (string, error) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", fmt.Errorf("error getting abs path of file: %w", err)
	}

	return absPath, nil
}

const subscriptionsFilepath = "internal/database/subscriptions.json"
const recipientsFilepath = "internal/database/recipients.json"
const sendHistoryFilepath = "internal/database/send_history.json"

func GetSubscriptions() ([]models.Subscription, error) {
	subscriptionsRaw, err := os.ReadFile(subscriptionsFilepath)
	if err != nil {
		return nil, fmt.Errorf("error reading subscriptions file: %w", err)
	}

	var subscriptions []subscription
	if err = json.Unmarshal(subscriptionsRaw, &subscriptions); err != nil {
		return nil, fmt.Errorf("error decoding subscriptions: %w", err)
	}

	recipientsRaw, err := os.ReadFile(recipientsFilepath)
	if err != nil {
		return nil, fmt.Errorf("error reading recipients file: %w", err)
	}

	var recipients []recipient
	if err = json.Unmarshal(recipientsRaw, &recipients); err != nil {
		return nil, fmt.Errorf("error unmarshalling recipients: %w", err)
	}

	recipientsMap := make(map[int]recipient, len(recipients))
	for _, r := range recipients {
		recipientsMap[r.ID] = r
	}

	hydratedSubscriptions := []models.Subscription{}
	for _, subscription := range subscriptions {
		for _, recipientID := range subscription.Recipients {
			recipient := recipientsMap[recipientID]
			hydratedSubscription := models.Subscription{
				Coach: models.Coach{
					ID:    recipientID,
					Name:  recipient.Name,
					Email: recipient.Email,
				},
				LeagueID: subscription.LeagueID,
				SeasonID: subscription.SeasonID,
				TeamName: subscription.TeamName,
			}
			hydratedSubscriptions = append(hydratedSubscriptions, hydratedSubscription)
		}
	}

	return hydratedSubscriptions, nil
}

type historyEntry struct {
	GameID      int `json:"game_id"`
	RecipientID int `json:"recipient_id"`
}

func UpdateSendHistory(gameID int, recipientID int) error {
	sendHistoryFile, err := os.OpenFile(sendHistoryFilepath, os.O_RDWR, 0644) // standard FileMode ? idk got from internet
	if err != nil {
		return err
	}
	defer func() {
		err := sendHistoryFile.Close()
		if err != nil {
			fmt.Println("error closing history file: %w", err)
		}
	}()

	var history []historyEntry

	if err := json.NewDecoder(sendHistoryFile).Decode(&history); err != nil {
		return fmt.Errorf("error decoding history file: %w", err)
	}

	history = append(history, historyEntry{GameID: gameID, RecipientID: recipientID})

	bytes, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("error marshalling history")
	}

	_, err = sendHistoryFile.WriteAt(bytes, 0)

	return err
}

func HasSentAlert(gameID int, recipientID int) (bool, error) {
	raw, err := os.ReadFile(sendHistoryFilepath)
	if err != nil {
		return false, fmt.Errorf("error reading history file: %w", err)
	}

	var history []historyEntry
	err = json.Unmarshal(raw, &history)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling history data: %w", err)
	}

	for _, h := range history {
		if h.GameID == gameID && h.RecipientID == recipientID {
			return true, nil
		}
	}

	return false, nil
}
