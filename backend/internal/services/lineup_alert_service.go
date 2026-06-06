package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

type LineupAlertConfig struct {
	Enabled       bool
	Interval      time.Duration
	WindowBefore  time.Duration
	WindowAfter   time.Duration
}

var lineupAlertCfg LineupAlertConfig

func ConfigureLineupAlert(cfg LineupAlertConfig) {
	lineupAlertCfg = cfg
}

func IsLineupAlertEnabled() bool {
	return lineupAlertCfg.Enabled
}

func NextLineupAlertInterval() time.Duration {
	return lineupAlertCfg.Interval
}

// ScanAndSendLineupAlerts scans matches within the alert window and sends
// lineup notifications to interested users when lineups become available.
func ScanAndSendLineupAlerts(ctx context.Context) error {
	if !lineupAlertCfg.Enabled {
		return nil
	}

	now := time.Now().UTC()
	startWindow := now.Add(-lineupAlertCfg.WindowAfter)
	endWindow := now.Add(lineupAlertCfg.WindowBefore)

	matches, err := repositories.ListMatchesForLineupAlert(startWindow, endWindow)
	if err != nil {
		return fmt.Errorf("query matches for lineup alert: %w", err)
	}

	for _, match := range matches {
		if err := sendLineupAlertForMatch(ctx, match); err != nil {
			log.Printf("Lineup alert failed for match %d: %v", match.ID, err)
		}
	}
	return nil
}

// sendLineupAlertForMatch checks if lineups are ready and sends notifications.
func sendLineupAlertForMatch(ctx context.Context, match models.Match) error {
	lineups, err := repositories.GetLineupsByMatch(match.ID)
	if err != nil {
		return err
	}

	if !isLineupAlertReady(lineups) {
		return nil
	}

	// Collect target users: favorite match + reminders
	userIDs, err := collectTargetUserIDs(match.ID)
	if err != nil {
		return err
	}
	if len(userIDs) == 0 {
		return nil
	}

	title := fmt.Sprintf("首发已公布：%s vs %s", match.HomeTeam.Name, match.AwayTeam.Name)
	content := "双方首发阵容已更新，点击查看阵型和首发名单。"

	// Dedup by UserMatchEventLog and send notifications in a transaction
	for _, userID := range userIDs {
		exists, err := repositories.HasUserMatchEvent(userID, match.ID, "lineup_available")
		if err != nil || exists {
			continue
		}

		notification := &models.Notification{
			UserID:     userID,
			Title:      title,
			Content:    content,
			Type:       "lineup",
			TargetType: "match",
			TargetID:   match.ID,
		}
		if err := repositories.CreateNotification(notification); err != nil {
			log.Printf("Create lineup notification failed for user %d match %d: %v", userID, match.ID, err)
			continue
		}
		_ = repositories.CreateUserMatchEvent(userID, match.ID, "lineup_available")
	}
	return nil
}

// isLineupAlertReady checks that both home and away have ≥ 11 starting players.
func isLineupAlertReady(lineups []models.MatchLineup) bool {
	var homeStarting, awayStarting int
	for _, lu := range lineups {
		count := 0
		for _, p := range lu.Players {
			if p.Role == "starting" {
				count++
			}
		}
		if lu.Side == "home" {
			homeStarting = count
		} else if lu.Side == "away" {
			awayStarting = count
		}
	}
	return homeStarting >= 11 && awayStarting >= 11
}

// collectTargetUserIDs gathers user IDs who favorited or set reminders for a match.
func collectTargetUserIDs(matchID uint) ([]uint, error) {
	seen := make(map[uint]bool)
	var allIDs []uint

	favIDs, err := repositories.ListUserIDsByFavoriteMatch(matchID)
	if err != nil {
		return nil, err
	}
	for _, id := range favIDs {
		if !seen[id] {
			seen[id] = true
			allIDs = append(allIDs, id)
		}
	}

	remindIDs, err := repositories.ListUserIDsByMatchReminders(matchID)
	if err != nil {
		return nil, err
	}
	for _, id := range remindIDs {
		if !seen[id] {
			seen[id] = true
			allIDs = append(allIDs, id)
		}
	}

	return allIDs, nil
}
