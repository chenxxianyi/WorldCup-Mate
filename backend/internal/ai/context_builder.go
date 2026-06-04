package ai

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

type ContextBuilder struct{}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{}
}

func (b *ContextBuilder) MatchContext(matchID uint, userID *uint) (string, *models.Match, error) {
	var match models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").Preload("Group").
		First(&match, matchID).Error
	if err != nil {
		return "", nil, err
	}

	lines := []string{
		fmt.Sprintf("Match ID: %d", match.ID),
		fmt.Sprintf("Match No: %d", match.MatchNo),
		fmt.Sprintf("Stage: %s", clean(match.Stage)),
		fmt.Sprintf("Status: %s", clean(match.Status)),
		fmt.Sprintf("Kickoff UTC: %s", match.KickoffTimeUTC.UTC().Format(time.RFC3339)),
		fmt.Sprintf("Home team: %s (%s)", teamName(match.HomeTeam), match.HomeTeam.FIFACode),
		fmt.Sprintf("Away team: %s (%s)", teamName(match.AwayTeam), match.AwayTeam.FIFACode),
		fmt.Sprintf("City: %s", nameOrEn(match.City.Name, match.City.NameEn)),
		fmt.Sprintf("Stadium: %s", nameOrEn(match.Stadium.Name, match.Stadium.NameEn)),
	}
	if match.Group != nil {
		lines = append(lines, fmt.Sprintf("Group: %s", clean(match.Group.Name)))
	}
	if match.HomeScore != nil && match.AwayScore != nil {
		lines = append(lines, fmt.Sprintf("Score: %d-%d", *match.HomeScore, *match.AwayScore))
	}
	if match.ImportanceLevel > 0 {
		lines = append(lines, fmt.Sprintf("Importance level: %d", match.ImportanceLevel))
	}
	if text := clean(match.RecommendTag); text != "" {
		lines = append(lines, "Recommend tag: "+text)
	}
	if text := clean(match.RecommendReason); text != "" {
		lines = append(lines, "Recommend reason: "+text)
	}

	if match.GroupID != nil {
		standings := b.groupStandingLines(*match.GroupID)
		if len(standings) > 0 {
			lines = append(lines, "Group standings:")
			lines = append(lines, standings...)
		} else {
			lines = append(lines, "Group standings: no standings data yet")
		}
	}
	if userID != nil {
		lines = append(lines, b.favoriteTeamLine(*userID))
	}
	return strings.Join(lines, "\n"), &match, nil
}

func (b *ContextBuilder) TodayMatchesContext(date, timezone string, userID *uint) (string, []models.Match, error) {
	start, err := localDayStart(date, timezone)
	if err != nil {
		return "", nil, err
	}
	end := start.AddDate(0, 0, 1)
	var matches []models.Match
	err = database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").Preload("Group").
		Where("kickoff_time_utc >= ? AND kickoff_time_utc < ?", start.UTC(), end.UTC()).
		Order("importance_level DESC, kickoff_time_utc ASC").
		Find(&matches).Error
	if err != nil {
		return "", nil, err
	}
	if len(matches) == 0 {
		err = database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").Preload("Group").
			Where("kickoff_time_utc >= ?", end.UTC()).
			Order("kickoff_time_utc ASC").
			Limit(5).
			Find(&matches).Error
		if err != nil {
			return "", nil, err
		}
	}

	lines := []string{
		fmt.Sprintf("Date: %s", date),
		fmt.Sprintf("Timezone: %s", timezone),
	}
	if userID != nil {
		lines = append(lines, b.favoriteTeamLine(*userID))
	}
	if len(matches) == 0 {
		lines = append(lines, "Matches: no matches found in database")
	} else {
		lines = append(lines, "Matches:")
		for _, m := range matches {
			lines = append(lines, matchLine(m))
		}
	}
	return strings.Join(lines, "\n"), matches, nil
}

func (b *ContextBuilder) GroupContext(groupID uint) (string, *models.Group, []models.GroupStanding, []models.Match, error) {
	var group models.Group
	if err := database.DB.First(&group, groupID).Error; err != nil {
		return "", nil, nil, nil, err
	}
	var standings []models.GroupStanding
	if err := database.DB.Preload("Team").Where("group_id = ?", groupID).Order("`rank` ASC").Find(&standings).Error; err != nil {
		return "", nil, nil, nil, err
	}
	var matches []models.Match
	if err := database.DB.Preload("HomeTeam").Preload("AwayTeam").
		Where("group_id = ?", groupID).Order("kickoff_time_utc ASC").Find(&matches).Error; err != nil {
		return "", nil, nil, nil, err
	}

	lines := []string{
		fmt.Sprintf("Group: %s", clean(group.Name)),
		"Qualification rule: top two teams qualify directly; best third-place teams may qualify according to tournament rules.",
	}
	if len(standings) == 0 {
		lines = append(lines, "Standings: no standings data yet")
	} else {
		lines = append(lines, "Standings:")
		for _, s := range standings {
			lines = append(lines, fmt.Sprintf("%d. %s %d pts, GD %d, status %s", s.Rank, teamName(s.Team), s.Points, s.GoalDifference, clean(s.QualificationStatus)))
		}
	}
	if len(matches) > 0 {
		lines = append(lines, "Group matches:")
		for _, m := range matches {
			lines = append(lines, matchLine(m))
		}
	}
	return strings.Join(lines, "\n"), &group, standings, matches, nil
}

func (b *ContextBuilder) ChatContext(contextType string, contextID uint, userID uint) string {
	switch contextType {
	case "match":
		ctx, _, err := b.MatchContext(contextID, &userID)
		if err == nil {
			return ctx
		}
	case "group":
		ctx, _, _, _, err := b.GroupContext(contextID)
		if err == nil {
			return ctx
		}
	case "team":
		var team models.Team
		if err := database.DB.Preload("Group").First(&team, contextID).Error; err == nil {
			return fmt.Sprintf("Team: %s (%s)\nGroup: %s", teamName(team), team.FIFACode, clean(team.Group.Name))
		}
	}
	return b.favoriteTeamLine(userID)
}

func (b *ContextBuilder) groupStandingLines(groupID uint) []string {
	var standings []models.GroupStanding
	if err := database.DB.Preload("Team").Where("group_id = ?", groupID).Order("`rank` ASC").Find(&standings).Error; err != nil {
		return nil
	}
	lines := make([]string, 0, len(standings))
	for _, s := range standings {
		lines = append(lines, fmt.Sprintf("- %d. %s: %d pts, GD %d, status %s", s.Rank, teamName(s.Team), s.Points, s.GoalDifference, clean(s.QualificationStatus)))
	}
	return lines
}

func (b *ContextBuilder) favoriteTeamLine(userID uint) string {
	var favs []models.UserFavoriteTeam
	if err := database.DB.Preload("Team").Where("user_id = ?", userID).Find(&favs).Error; err != nil || len(favs) == 0 {
		return "Favorite teams: none"
	}
	names := make([]string, 0, len(favs))
	for _, fav := range favs {
		names = append(names, teamName(fav.Team))
	}
	return "Favorite teams: " + strings.Join(names, ", ")
}

func localDayStart(date, timezone string) (time.Time, error) {
	if strings.TrimSpace(timezone) == "" {
		timezone = "Asia/Shanghai"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.FixedZone(timezone, 8*3600)
	}
	if strings.TrimSpace(date) == "" {
		now := time.Now().In(loc)
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc), nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func matchLine(m models.Match) string {
	return fmt.Sprintf("- Match ID %d, Match No %d: %s vs %s at %s UTC, status %s, importance %d", m.ID, m.MatchNo, teamName(m.HomeTeam), teamName(m.AwayTeam), m.KickoffTimeUTC.UTC().Format(time.RFC3339), clean(m.Status), m.ImportanceLevel)
}

func teamName(t models.Team) string {
	return nameOrEn(t.Name, t.NameEn)
}

func nameOrEn(name, en string) string {
	if isCleanText(name) {
		return strings.TrimSpace(name)
	}
	if strings.TrimSpace(en) != "" {
		return strings.TrimSpace(en)
	}
	return strings.TrimSpace(name)
}

func clean(s string) string {
	if isCleanText(s) {
		return strings.TrimSpace(s)
	}
	return ""
}

func isCleanText(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if !utf8.ValidString(s) {
		return false
	}
	return !strings.Contains(s, "�") && !strings.Contains(s, "锟") && !strings.Contains(s, "涓") && !strings.Contains(s, "鍦")
}
