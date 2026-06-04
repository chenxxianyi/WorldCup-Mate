package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func AIMatchInsight(c *gin.Context) {
	var req services.MatchInsightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}
	res, err := services.GenerateMatchInsight(c.Request.Context(), req, optionalUserID(c), c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}

func AITodayRecommendations(c *gin.Context) {
	var req services.TodayRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}
	res, err := services.GenerateTodayRecommendations(c.Request.Context(), req, optionalUserID(c), c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}

func AIGroupAnalysis(c *gin.Context) {
	var req services.GroupAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}
	res, err := services.GenerateGroupAnalysis(c.Request.Context(), req, optionalUserID(c), c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}

func AIExplain(c *gin.Context) {
	var req services.ExplainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}
	res, err := services.ExplainFootball(c.Request.Context(), req, optionalUserID(c), c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}

func AIShareCopy(c *gin.Context) {
	var req services.ShareCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}
	res, err := services.GenerateShareCopy(c.Request.Context(), req, optionalUserID(c), c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}

func AIChat(c *gin.Context) {
	userID, ok := requiredUserID(c)
	if !ok {
		return
	}
	var req services.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}
	res, err := services.Chat(c.Request.Context(), req, userID, c.ClientIP())
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, res)
}

func AIChatStream(c *gin.Context) {
	userID, ok := requiredUserID(c)
	if !ok {
		return
	}
	var req services.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "请求参数不正确")
		return
	}

	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flush := func(event services.AIChatStreamEvent) error {
		body, err := json.Marshal(event)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", body); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}

	if _, err := services.ChatStream(c.Request.Context(), req, userID, c.ClientIP(), flush); err != nil {
		_ = flush(services.AIChatStreamEvent{Type: "error", Delta: err.Error()})
	}
}

func AIListConversations(c *gin.Context) {
	userID, ok := requiredUserID(c)
	if !ok {
		return
	}
	res, err := services.ListAIConversations(userID)
	if err != nil {
		utils.Error(c, 500, "读取会话失败")
		return
	}
	utils.Success(c, res)
}

func AIGetConversation(c *gin.Context) {
	userID, ok := requiredUserID(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Error(c, 400, "invalid conversation id")
		return
	}
	res, err := services.GetAIConversation(userID, uint(id))
	if err != nil {
		utils.Error(c, 404, "会话不存在")
		return
	}
	utils.Success(c, res)
}

func AIDeleteConversation(c *gin.Context) {
	userID, ok := requiredUserID(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.Error(c, 400, "invalid conversation id")
		return
	}
	if err := services.DeleteAIConversation(userID, uint(id)); err != nil {
		utils.Error(c, 404, "会话不存在")
		return
	}
	utils.Success(c, nil)
}

func optionalUserID(c *gin.Context) *uint {
	if value, ok := c.Get("user_id"); ok {
		if id, ok := value.(uint); ok && id > 0 {
			return &id
		}
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return nil
	}
	claims, err := utils.ParseToken(tokenStr)
	if err != nil || claims.UserID == 0 {
		return nil
	}
	id := claims.UserID
	return &id
}

func requiredUserID(c *gin.Context) (uint, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, 401, "请先登录")
		return 0, false
	}
	userID, ok := value.(uint)
	if !ok || userID == 0 {
		utils.Error(c, 401, "请先登录")
		return 0, false
	}
	return userID, true
}
