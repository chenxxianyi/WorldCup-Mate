package handlers

import (
	"strings"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

// PublicFeaturedInput / admin upsert payload.
type FeaturedInput struct {
	MatchID     *uint  `json:"match_id"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	StageLabel  string `json:"stage_label"`
	Enabled     *bool  `json:"enabled"`
}

// PublicFeatured returns every enabled hero config, keyed by competition
// code (frontend merges it over its static theme copy).
func PublicFeatured(c *gin.Context) {
	configs, err := repositories.GetFeaturedConfigs()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	result := map[string]*models.FeaturedConfig{}
	for _, cfg := range configs {
		if cfg.Enabled {
			cfg := cfg
			result[cfg.Competition.Code] = &cfg
		}
	}
	utils.Success(c, result)
}

// AdminListFeatured returns all configs with competition metadata.
func AdminListFeatured(c *gin.Context) {
	configs, err := repositories.GetFeaturedConfigs()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, configs)
}

// AdminUpsertFeatured creates/updates the hero config for a competition.
func AdminUpsertFeatured(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	competition, err := repositories.GetCompetitionByCode(code)
	if err != nil {
		utils.Error(c, 404, "competition not found")
		return
	}
	var input FeaturedInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	// The pinned match must belong to this competition.
	if input.MatchID != nil {
		match, err := repositories.GetMatchByID(*input.MatchID)
		if err != nil {
			utils.Error(c, 400, "match not found")
			return
		}
		if match.CompetitionID == nil || *match.CompetitionID != competition.ID {
			utils.Error(c, 400, "match does not belong to this competition")
			return
		}
	}
	cfg := &models.FeaturedConfig{
		MatchID:     input.MatchID,
		Tagline:     input.Tagline,
		Description: input.Description,
		StageLabel:  input.StageLabel,
		Enabled:     true,
	}
	if input.Enabled != nil {
		cfg.Enabled = *input.Enabled
	}
	if err := repositories.UpsertFeaturedConfig(competition.ID, cfg); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "featured", code, "upsert", nil, gin.H{
		"match_id":     cfg.MatchID,
		"tagline":      cfg.Tagline,
		"description":  cfg.Description,
		"stage_label":  cfg.StageLabel,
		"enabled":      cfg.Enabled,
	})
	utils.Success(c, cfg)
}

// AdminFeaturedMatches lists a competition's matches for the focus-match
// picker: upcoming (kickoff ascending) first, then finished (recent first),
// so the "next focus match" is never pushed out of the top 50.
func AdminFeaturedMatches(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	competition, err := repositories.GetCompetitionByCode(code)
	if err != nil {
		utils.Error(c, 404, "competition not found")
		return
	}
	type pick struct {
		ID      uint   `json:"id"`
		Home    string `json:"home"`
		Away    string `json:"away"`
		Kickoff string `json:"kickoff_time_utc"`
		Status  string `json:"status"`
	}
	toPicks := func(ms []models.Match) []pick {
		result := make([]pick, 0, len(ms))
		for _, m := range ms {
			result = append(result, pick{
				ID:      m.ID,
				Home:    m.HomeTeam.Name,
				Away:    m.AwayTeam.Name,
				Kickoff: m.KickoffTimeUTC.Format("2006-01-02T15:04:05Z07:00"),
				Status:  m.Status,
			})
		}
		return result
	}
	upcoming, _, err := repositories.ListMatches(repositories.MatchFilter{
		CompetitionID: competition.ID,
		Season:        competition.Season,
		Status:        "scheduled",
		Page:          1,
		PageSize:      50,
	})
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	finished, _, err := repositories.ListMatches(repositories.MatchFilter{
		CompetitionID: competition.ID,
		Season:        competition.Season,
		Status:        "finished",
		Page:          1,
		PageSize:      50,
	})
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	live, _, err := repositories.ListMatches(repositories.MatchFilter{
		CompetitionID: competition.ID,
		Season:        competition.Season,
		Status:        "live",
		Page:          1,
		PageSize:      20,
	})
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	// Finished matches come back in kickoff-ascending order; reverse to
	// show the most recent first.
	for i, j := 0, len(finished)-1; i < j; i, j = i+1, j-1 {
		finished[i], finished[j] = finished[j], finished[i]
	}
	picks := append(toPicks(upcoming), toPicks(live)...)
	picks = append(picks, toPicks(finished)...)
	utils.Success(c, picks)
}
