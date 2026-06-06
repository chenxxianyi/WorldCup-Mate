package handlers

import (
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetPostMatchSummary handles GET /api/matches/:id/post-match-summary
// Returns existing summary from cache/DB if available, or empty state.
func GetPostMatchSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Error(c, 400, "invalid match id")
		return
	}

	// Try to fetch with ForceRefresh=false, but don't generate if not exists.
	// GeneratePostMatchSummary with ForceRefresh=false will return cached if exists,
	// or error if not. We treat "no summary" as a valid empty response.
	res, err := services.GeneratePostMatchSummary(c.Request.Context(), services.PostMatchSummaryRequest{
		MatchID:      uint(id),
		ForceRefresh: false,
	}, optionalUserID(c), c.ClientIP())
	if err != nil {
		// Return empty state - no summary available yet
		utils.Success(c, gin.H{
			"summary":       "",
			"score_line":    "",
			"key_takeaways": []string{},
			"data_note":     "",
			"generated_at":  "",
		})
		return
	}
	utils.Success(c, res)
}

// GeneratePostMatchSummary handles POST /api/matches/:id/post-match-summary/generate
// Actively generates a new post-match summary.
func GeneratePostMatchSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Error(c, 400, "invalid match id")
		return
	}

	var body struct {
		ForceRefresh bool `json:"force_refresh"`
	}
	_ = c.ShouldBindJSON(&body)

	res, err := services.GeneratePostMatchSummary(c.Request.Context(), services.PostMatchSummaryRequest{
		MatchID:      uint(id),
		ForceRefresh: body.ForceRefresh,
	}, optionalUserID(c), c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}
