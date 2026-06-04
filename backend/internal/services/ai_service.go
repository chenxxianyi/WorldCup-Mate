package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"worldcup-mate/internal/ai"
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AIServiceConfig struct {
	Provider       string
	BaseURL        string
	APIKey         string
	Model          string
	TimeoutSeconds int
	DailyLimitUser int
	Temperature    float64
	MaxTokens      int
	CacheEnabled   bool
}

type MatchInsightRequest struct {
	MatchID      uint `json:"match_id" binding:"required"`
	ForceRefresh bool `json:"force_refresh"`
}

type MatchInsightResponse struct {
	Summary             string    `json:"summary"`
	WatchRating         int       `json:"watch_rating"`
	Reasons             []string  `json:"reasons"`
	TeamComparison      []string  `json:"team_comparison"`
	BeginnerTips        []string  `json:"beginner_tips"`
	QualificationImpact string    `json:"qualification_impact"`
	ShouldStayUp        string    `json:"should_stay_up"`
	SuitableFor         []string  `json:"suitable_for"`
	GeneratedAt         time.Time `json:"generated_at"`
}

type TodayRecommendationRequest struct {
	Date         string `json:"date"`
	Timezone     string `json:"timezone"`
	Limit        int    `json:"limit"`
	ForceRefresh bool   `json:"force_refresh"`
}

type TodayRecommendedMatch struct {
	MatchID     uint   `json:"match_id"`
	Title       string `json:"title"`
	KickoffTime string `json:"kickoff_time,omitempty"`
	Reason      string `json:"reason"`
	Rating      int    `json:"rating,omitempty"`
}

type TodayRecommendationResponse struct {
	Date            string                  `json:"date"`
	Timezone        string                  `json:"timezone"`
	Recommendations []TodayRecommendedMatch `json:"recommendations"`
	OnlyOneMatch    *TodayRecommendedMatch  `json:"only_one_match"`
	Note            string                  `json:"note,omitempty"`
}

type GroupAnalysisRequest struct {
	GroupID      uint `json:"group_id" binding:"required"`
	ForceRefresh bool `json:"force_refresh"`
}

type GroupAnalysisTeam struct {
	TeamID   uint   `json:"team_id,omitempty"`
	TeamName string `json:"team_name"`
	Status   string `json:"status,omitempty"`
	Note     string `json:"note,omitempty"`
}

type GroupAnalysisResponse struct {
	Summary            string              `json:"summary"`
	KeyPoints          []string            `json:"key_points"`
	QualificationRules string              `json:"qualification_rules"`
	Teams              []GroupAnalysisTeam `json:"teams"`
	DataNote           string              `json:"data_note,omitempty"`
	GeneratedAt        time.Time           `json:"generated_at"`
}

type ExplainRequest struct {
	Question    string `json:"question" binding:"required"`
	ContextType string `json:"context_type"`
	ContextID   uint   `json:"context_id"`
}

type ExplainResponse struct {
	Answer    string   `json:"answer"`
	KeyPoints []string `json:"key_points,omitempty"`
}

type ShareCopyRequest struct {
	MatchID  uint   `json:"match_id" binding:"required"`
	Platform string `json:"platform" binding:"required"`
	Tone     string `json:"tone" binding:"required"`
	Length   string `json:"length" binding:"required"`
}

type ShareCopyResponse struct {
	Title   string   `json:"title,omitempty"`
	Content string   `json:"content"`
	Tips    []string `json:"tips,omitempty"`
}

type AIChatRequest struct {
	ConversationID *uint  `json:"conversation_id"`
	Message        string `json:"message" binding:"required"`
	ContextType    string `json:"context_type"`
	ContextID      uint   `json:"context_id"`
}

type AIChatMessageResponse struct {
	ID        uint      `json:"id,omitempty"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AIChatResponse struct {
	ConversationID uint                  `json:"conversation_id"`
	Message        AIChatMessageResponse `json:"message"`
}

type AIConversationResponse struct {
	ID          uint                    `json:"id"`
	Title       string                  `json:"title"`
	ContextType string                  `json:"context_type"`
	ContextID   uint                    `json:"context_id,omitempty"`
	LastMessage string                  `json:"last_message,omitempty"`
	Messages    []AIChatMessageResponse `json:"messages,omitempty"`
	CreatedAt   time.Time               `json:"created_at,omitempty"`
	UpdatedAt   time.Time               `json:"updated_at,omitempty"`
}

var aiSvc *AIService

type AIService struct {
	cfg      AIServiceConfig
	provider ai.Provider
	builder  *ai.ContextBuilder
}

func ConfigureAI(cfg AIServiceConfig) error {
	provider, err := ai.NewProvider(ai.ProviderConfig{
		Provider:       cfg.Provider,
		BaseURL:        cfg.BaseURL,
		APIKey:         cfg.APIKey,
		Model:          cfg.Model,
		TimeoutSeconds: cfg.TimeoutSeconds,
	})
	if err != nil {
		return err
	}
	if cfg.DailyLimitUser <= 0 {
		cfg.DailyLimitUser = 50
	}
	if cfg.MaxTokens <= 0 {
		cfg.MaxTokens = 1200
	}
	aiSvc = &AIService{
		cfg:      cfg,
		provider: provider,
		builder:  ai.NewContextBuilder(),
	}
	return nil
}

func GenerateMatchInsight(ctx context.Context, req MatchInsightRequest, userID *uint, ip string) (*MatchInsightResponse, error) {
	return currentAI().GenerateMatchInsight(ctx, req, userID, ip)
}

func GenerateTodayRecommendations(ctx context.Context, req TodayRecommendationRequest, userID *uint, ip string) (*TodayRecommendationResponse, error) {
	return currentAI().GenerateTodayRecommendations(ctx, req, userID, ip)
}

func GenerateGroupAnalysis(ctx context.Context, req GroupAnalysisRequest, userID *uint, ip string) (*GroupAnalysisResponse, error) {
	return currentAI().GenerateGroupAnalysis(ctx, req, userID, ip)
}

func ExplainFootball(ctx context.Context, req ExplainRequest, userID *uint, ip string) (*ExplainResponse, error) {
	return currentAI().ExplainFootball(ctx, req, userID, ip)
}

func GenerateShareCopy(ctx context.Context, req ShareCopyRequest, userID *uint, ip string) (*ShareCopyResponse, error) {
	return currentAI().GenerateShareCopy(ctx, req, userID, ip)
}

func Chat(ctx context.Context, req AIChatRequest, userID uint, ip string) (*AIChatResponse, error) {
	return currentAI().Chat(ctx, req, userID, ip)
}

func ListAIConversations(userID uint) ([]AIConversationResponse, error) {
	convs, err := repositories.ListConversations(userID)
	if err != nil {
		return nil, err
	}
	out := make([]AIConversationResponse, 0, len(convs))
	for _, conv := range convs {
		out = append(out, conversationResponse(conv, nil))
	}
	return out, nil
}

func GetAIConversation(userID, conversationID uint) (*AIConversationResponse, error) {
	conv, messages, err := repositories.GetConversationWithMessages(userID, conversationID)
	if err != nil {
		return nil, err
	}
	res := conversationResponse(*conv, messages)
	return &res, nil
}

func DeleteAIConversation(userID, conversationID uint) error {
	return repositories.DeleteConversation(userID, conversationID)
}

func currentAI() *AIService {
	if aiSvc == nil {
		_ = ConfigureAI(AIServiceConfig{Provider: "openai", Model: "gpt-4o-mini", DailyLimitUser: 50, MaxTokens: 1200, CacheEnabled: true})
	}
	return aiSvc
}

func (s *AIService) GenerateMatchInsight(ctx context.Context, req MatchInsightRequest, userID *uint, ip string) (*MatchInsightResponse, error) {
	if req.MatchID == 0 {
		return nil, fmt.Errorf("请选择比赛")
	}
	cacheKey := s.cacheKey("match_insight", fmt.Sprintf("%d", req.MatchID))
	var cached MatchInsightResponse
	if s.getCached(ctx, cacheKey, req.ForceRefresh, &cached) {
		return &cached, nil
	}

	matchCtx, _, err := s.builder.MatchContext(req.MatchID, userID)
	if err != nil {
		s.logUsage(userID, ip, "match_insight", "failed", err, nil, 0)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("比赛不存在")
		}
		return nil, fmt.Errorf("读取比赛信息失败")
	}
	if err := s.checkLimit(ctx, userID); err != nil {
		return nil, err
	}

	prompt := "TASK:match_insight\nReturn JSON with summary, watch_rating, reasons, team_comparison, beginner_tips, qualification_impact, should_stay_up, suitable_for.\nContext:\n" + matchCtx
	raw, providerRes, latency, err := s.callProvider(ctx, prompt, nil)
	if err != nil {
		s.logUsage(userID, ip, "match_insight", "failed", err, providerRes, latency)
		return nil, aiUserError(err)
	}
	if err := ai.ValidateOutput(raw); err != nil {
		s.logUsage(userID, ip, "match_insight", "failed", err, providerRes, latency)
		return nil, fmt.Errorf("AI 输出未通过安全检查")
	}

	res, jsonText, err := ai.DecodeJSON(raw, MatchInsightResponse{})
	if err != nil || strings.TrimSpace(res.Summary) == "" {
		s.logUsage(userID, ip, "match_insight", "failed", fmt.Errorf("invalid AI response format"), providerRes, latency)
		return nil, fmt.Errorf("AI 返回内容格式不正确，请稍后重试")
	}
	res.WatchRating = ai.ClampInt(res.WatchRating, 1, 5)
	res.Reasons = nonNilStrings(res.Reasons)
	res.TeamComparison = nonNilStrings(res.TeamComparison)
	res.BeginnerTips = nonNilStrings(res.BeginnerTips)
	res.SuitableFor = nonNilStrings(res.SuitableFor)
	if res.GeneratedAt.IsZero() {
		res.GeneratedAt = time.Now().UTC()
	}
	s.saveGenerated(userID, "match_insight", "match", req.MatchID, cacheKey, jsonText, raw, providerRes)
	s.setCached(ctx, cacheKey, res, 2*time.Hour)
	s.logUsage(userID, ip, "match_insight", "success", nil, providerRes, latency)
	return &res, nil
}

func (s *AIService) GenerateTodayRecommendations(ctx context.Context, req TodayRecommendationRequest, userID *uint, ip string) (*TodayRecommendationResponse, error) {
	if req.Timezone == "" {
		req.Timezone = "Asia/Shanghai"
	}
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}
	if req.Limit <= 0 || req.Limit > 10 {
		req.Limit = 3
	}
	cacheKey := s.cacheKey("today_recommendations", req.Date, req.Timezone, userKey(userID))
	var cached TodayRecommendationResponse
	if s.getCached(ctx, cacheKey, req.ForceRefresh, &cached) {
		return &cached, nil
	}

	todayCtx, matches, err := s.builder.TodayMatchesContext(req.Date, req.Timezone, userID)
	if err != nil {
		s.logUsage(userID, ip, "today_recommendations", "failed", err, nil, 0)
		return nil, fmt.Errorf("读取赛程失败")
	}
	if err := s.checkLimit(ctx, userID); err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf("TASK:today_recommendations\nReturn JSON with date, timezone, recommendations, only_one_match, note. Limit recommendations to %d. Every recommendation.match_id must be one of the Match ID values in Context.\nContext:\n%s", req.Limit, todayCtx)
	raw, providerRes, latency, err := s.callProvider(ctx, prompt, nil)
	if err != nil {
		s.logUsage(userID, ip, "today_recommendations", "failed", err, providerRes, latency)
		return nil, aiUserError(err)
	}
	if err := ai.ValidateOutput(raw); err != nil {
		s.logUsage(userID, ip, "today_recommendations", "failed", err, providerRes, latency)
		return nil, fmt.Errorf("AI 输出未通过安全检查")
	}
	res, jsonText, err := ai.DecodeJSON(raw, TodayRecommendationResponse{})
	if err != nil {
		s.logUsage(userID, ip, "today_recommendations", "failed", fmt.Errorf("invalid AI response format"), providerRes, latency)
		return nil, fmt.Errorf("AI 返回内容格式不正确，请稍后重试")
	}
	res.Date = req.Date
	res.Timezone = req.Timezone
	res.Recommendations = validateRecommendations(res.Recommendations, matches, req.Limit)
	if len(matches) > 0 && len(res.Recommendations) == 0 {
		s.logUsage(userID, ip, "today_recommendations", "failed", fmt.Errorf("AI response has no valid match ids"), providerRes, latency)
		return nil, fmt.Errorf("AI 返回内容缺少有效比赛，请稍后重试")
	}
	if len(res.Recommendations) == 1 {
		res.OnlyOneMatch = &res.Recommendations[0]
	}
	s.saveGenerated(userID, "today_recommendations", "date", 0, cacheKey, jsonText, raw, providerRes)
	s.setCached(ctx, cacheKey, res, 30*time.Minute)
	s.logUsage(userID, ip, "today_recommendations", "success", nil, providerRes, latency)
	return &res, nil
}

func (s *AIService) GenerateGroupAnalysis(ctx context.Context, req GroupAnalysisRequest, userID *uint, ip string) (*GroupAnalysisResponse, error) {
	if req.GroupID == 0 {
		return nil, fmt.Errorf("请选择小组")
	}
	cacheKey := s.cacheKey("group_analysis", fmt.Sprintf("%d", req.GroupID))
	var cached GroupAnalysisResponse
	if s.getCached(ctx, cacheKey, req.ForceRefresh, &cached) {
		return &cached, nil
	}
	groupCtx, _, standings, _, err := s.builder.GroupContext(req.GroupID)
	if err != nil {
		s.logUsage(userID, ip, "group_analysis", "failed", err, nil, 0)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("小组不存在")
		}
		return nil, fmt.Errorf("读取小组信息失败")
	}
	if err := s.checkLimit(ctx, userID); err != nil {
		return nil, err
	}

	prompt := "TASK:group_analysis\nReturn JSON with summary, key_points, qualification_rules, teams, data_note. Do not invent probabilities.\nContext:\n" + groupCtx
	raw, providerRes, latency, err := s.callProvider(ctx, prompt, nil)
	if err != nil {
		s.logUsage(userID, ip, "group_analysis", "failed", err, providerRes, latency)
		return nil, aiUserError(err)
	}
	if err := ai.ValidateOutput(raw); err != nil {
		s.logUsage(userID, ip, "group_analysis", "failed", err, providerRes, latency)
		return nil, fmt.Errorf("AI 输出未通过安全检查")
	}
	res, jsonText, err := ai.DecodeJSON(raw, GroupAnalysisResponse{})
	if err != nil || strings.TrimSpace(res.Summary) == "" {
		s.logUsage(userID, ip, "group_analysis", "failed", fmt.Errorf("invalid AI response format"), providerRes, latency)
		return nil, fmt.Errorf("AI 返回内容格式不正确，请稍后重试")
	}
	if len(res.Teams) == 0 {
		res.Teams = buildGroupTeams(standings)
	}
	res.KeyPoints = nonNilStrings(res.KeyPoints)
	if res.GeneratedAt.IsZero() {
		res.GeneratedAt = time.Now().UTC()
	}
	s.saveGenerated(userID, "group_analysis", "group", req.GroupID, cacheKey, jsonText, raw, providerRes)
	s.setCached(ctx, cacheKey, res, 45*time.Minute)
	s.logUsage(userID, ip, "group_analysis", "success", nil, providerRes, latency)
	return &res, nil
}

func (s *AIService) ExplainFootball(ctx context.Context, req ExplainRequest, userID *uint, ip string) (*ExplainResponse, error) {
	req.Question = strings.TrimSpace(req.Question)
	if req.Question == "" {
		return nil, fmt.Errorf("请输入问题")
	}
	if err := s.checkLimit(ctx, userID); err != nil {
		return nil, err
	}
	contextText := ""
	if userID != nil {
		contextText = s.builder.ChatContext(defaultContext(req.ContextType), req.ContextID, *userID)
	}
	prompt := "TASK:explain\nReturn JSON with answer and key_points. Explain football rules in beginner-friendly Chinese.\nQuestion: " + req.Question + "\nContext:\n" + contextText
	raw, providerRes, latency, err := s.callProvider(ctx, prompt, nil)
	if err != nil {
		s.logUsage(userID, ip, "explain", "failed", err, providerRes, latency)
		return nil, aiUserError(err)
	}
	if err := ai.ValidateOutput(raw); err != nil {
		s.logUsage(userID, ip, "explain", "failed", err, providerRes, latency)
		return nil, fmt.Errorf("AI 输出未通过安全检查")
	}
	fallback := ExplainResponse{Answer: raw, KeyPoints: []string{}}
	res, _, _ := ai.DecodeJSON(raw, fallback)
	if res.Answer == "" {
		res.Answer = raw
	}
	s.logUsage(userID, ip, "explain", "success", nil, providerRes, latency)
	return &res, nil
}

func (s *AIService) GenerateShareCopy(ctx context.Context, req ShareCopyRequest, userID *uint, ip string) (*ShareCopyResponse, error) {
	if req.MatchID == 0 {
		return nil, fmt.Errorf("请选择比赛")
	}
	matchCtx, _, err := s.builder.MatchContext(req.MatchID, userID)
	if err != nil {
		s.logUsage(userID, ip, "share_copy", "failed", err, nil, 0)
		return nil, fmt.Errorf("比赛不存在")
	}
	if err := s.checkLimit(ctx, userID); err != nil {
		return nil, err
	}
	prompt := fmt.Sprintf("TASK:share_copy\nReturn JSON with title, content, tips. Platform=%s Tone=%s Length=%s. Do not mention betting or invented facts.\nContext:\n%s", req.Platform, req.Tone, req.Length, matchCtx)
	raw, providerRes, latency, err := s.callProvider(ctx, prompt, nil)
	if err != nil {
		s.logUsage(userID, ip, "share_copy", "failed", err, providerRes, latency)
		return nil, aiUserError(err)
	}
	if err := ai.ValidateOutput(raw); err != nil {
		s.logUsage(userID, ip, "share_copy", "failed", err, providerRes, latency)
		return nil, fmt.Errorf("AI 输出未通过安全检查")
	}
	res, jsonText, err := ai.DecodeJSON(raw, ShareCopyResponse{})
	if err != nil || strings.TrimSpace(res.Content) == "" {
		s.logUsage(userID, ip, "share_copy", "failed", fmt.Errorf("invalid AI response format"), providerRes, latency)
		return nil, fmt.Errorf("AI 返回内容格式不正确，请稍后重试")
	}
	cacheKey := s.cacheKey("share_copy", fmt.Sprintf("%d", req.MatchID), req.Platform, req.Tone, req.Length, userKey(userID))
	s.saveGenerated(userID, "share_copy", "match", req.MatchID, cacheKey, jsonText, raw, providerRes)
	s.logUsage(userID, ip, "share_copy", "success", nil, providerRes, latency)
	return &res, nil
}

func (s *AIService) Chat(ctx context.Context, req AIChatRequest, userID uint, ip string) (*AIChatResponse, error) {
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		return nil, fmt.Errorf("请输入消息")
	}
	if err := s.checkLimit(ctx, &userID); err != nil {
		return nil, err
	}
	contextType := defaultContext(req.ContextType)
	var conv *models.AIConversation
	var err error
	if req.ConversationID != nil && *req.ConversationID > 0 {
		conv, err = repositories.GetConversation(userID, *req.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("会话不存在")
		}
	} else {
		conv = &models.AIConversation{
			UserID:      userID,
			Title:       conversationTitle(req.Message),
			ContextType: contextType,
			ContextID:   req.ContextID,
		}
		if err := repositories.CreateConversation(conv); err != nil {
			return nil, fmt.Errorf("创建会话失败")
		}
	}

	userMsg := &models.AIMessage{ConversationID: conv.ID, UserID: userID, Role: "user", Content: req.Message}
	if err := repositories.SaveMessage(userMsg); err != nil {
		return nil, fmt.Errorf("保存消息失败")
	}

	recent, _ := repositories.ListRecentMessages(userID, conv.ID, 8)
	messages := make([]ai.Message, 0, len(recent))
	for _, msg := range recent {
		messages = append(messages, ai.Message{Role: msg.Role, Content: msg.Content})
	}
	factContext := s.builder.ChatContext(contextType, req.ContextID, userID)
	prompt := "TASK:chat\nUse only necessary facts. Answer in concise Chinese.\nFacts:\n" + factContext
	raw, providerRes, latency, err := s.callProvider(ctx, prompt, messages)
	if err != nil {
		s.logUsage(&userID, ip, "chat", "failed", err, providerRes, latency)
		return nil, aiUserError(err)
	}
	if err := ai.ValidateOutput(raw); err != nil {
		s.logUsage(&userID, ip, "chat", "failed", err, providerRes, latency)
		return nil, fmt.Errorf("AI 输出未通过安全检查")
	}
	assistantMsg := &models.AIMessage{
		ConversationID:   conv.ID,
		UserID:           userID,
		Role:             "assistant",
		Content:          raw,
		Provider:         s.provider.Name(),
		Model:            providerModel(providerRes, s.cfg.Model),
		PromptTokens:     providerTokens(providerRes, "prompt"),
		CompletionTokens: providerTokens(providerRes, "completion"),
		TotalTokens:      providerTokens(providerRes, "total"),
	}
	if err := repositories.SaveMessage(assistantMsg); err != nil {
		return nil, fmt.Errorf("保存回复失败")
	}
	conv.LastMessage = raw
	_ = repositories.UpdateConversation(conv)
	s.logUsage(&userID, ip, "chat", "success", nil, providerRes, latency)
	return &AIChatResponse{
		ConversationID: conv.ID,
		Message:        messageResponse(*assistantMsg),
	}, nil
}

func (s *AIService) callProvider(ctx context.Context, prompt string, messages []ai.Message) (string, *ai.ChatResponse, int64, error) {
	start := time.Now()
	timeout := time.Duration(s.cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := s.provider.Chat(callCtx, ai.ChatRequest{
		SystemPrompt: ai.BuildSystemPrompt(),
		UserPrompt:   prompt,
		Messages:     messages,
		Model:        s.cfg.Model,
		Temperature:  s.cfg.Temperature,
		MaxTokens:    s.cfg.MaxTokens,
	})
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return "", res, latency, err
	}
	return res.Content, res, latency, nil
}

func (s *AIService) checkLimit(ctx context.Context, userID *uint) error {
	if userID == nil || *userID == 0 || database.RDB == nil {
		return nil
	}
	key := fmt.Sprintf("ai:limit:user:%d:%s", *userID, time.Now().Format("2006-01-02"))
	n, err := database.RDB.Incr(ctx, key).Result()
	if err != nil {
		return nil
	}
	if n == 1 {
		_ = database.RDB.Expire(ctx, key, 26*time.Hour).Err()
	}
	if n > int64(s.cfg.DailyLimitUser) {
		return fmt.Errorf("今天的 AI 使用次数已达上限，明天再来")
	}
	return nil
}

func (s *AIService) getCached(ctx context.Context, key string, force bool, dest interface{}) bool {
	if force || !s.cfg.CacheEnabled {
		return false
	}
	if database.RDB != nil {
		raw, err := database.RDB.Get(ctx, key).Result()
		if err == nil && json.Unmarshal([]byte(raw), dest) == nil {
			return true
		}
		if err != nil && !errors.Is(err, redis.Nil) {
			return false
		}
	}
	content, err := repositories.GetGeneratedContentByCacheKey(key)
	if err == nil && content.ContentJSON != "" && json.Unmarshal([]byte(content.ContentJSON), dest) == nil {
		return true
	}
	return false
}

func (s *AIService) setCached(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	if !s.cfg.CacheEnabled || database.RDB == nil {
		return
	}
	body, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = database.RDB.Set(ctx, key, body, ttl).Err()
}

func (s *AIService) cacheKey(parts ...string) string {
	keyParts := []string{"ai", sanitizeCachePart(s.provider.Name())}
	keyParts = append(keyParts, parts...)
	for i := range keyParts {
		keyParts[i] = sanitizeCachePart(keyParts[i])
	}
	return strings.Join(keyParts, ":")
}

func (s *AIService) saveGenerated(userID *uint, typ, targetType string, targetID uint, cacheKey, jsonText, markdown string, providerRes *ai.ChatResponse) {
	content := &models.AIGeneratedContent{
		UserID:          userID,
		Type:            typ,
		TargetType:      targetType,
		TargetID:        targetID,
		ContentJSON:     jsonText,
		ContentMarkdown: markdown,
		Provider:        s.provider.Name(),
		Model:           providerModel(providerRes, s.cfg.Model),
		CacheKey:        cacheKey,
	}
	_ = repositories.SaveGeneratedContent(content)
}

func (s *AIService) logUsage(userID *uint, ip, endpoint, status string, callErr error, providerRes *ai.ChatResponse, latency int64) {
	log := &models.AIUsageLog{
		UserID:           userID,
		IP:               ip,
		Endpoint:         endpoint,
		Provider:         s.provider.Name(),
		Model:            providerModel(providerRes, s.cfg.Model),
		PromptTokens:     providerTokens(providerRes, "prompt"),
		CompletionTokens: providerTokens(providerRes, "completion"),
		TotalTokens:      providerTokens(providerRes, "total"),
		Status:           status,
		ErrorMessage:     ai.SanitizeError(callErr),
		LatencyMS:        latency,
	}
	_ = repositories.SaveUsageLog(log)
}

func providerModel(res *ai.ChatResponse, fallback string) string {
	if res != nil && res.Model != "" {
		return res.Model
	}
	return fallback
}

func providerTokens(res *ai.ChatResponse, kind string) int {
	if res == nil {
		return 0
	}
	switch kind {
	case "prompt":
		return res.PromptTokens
	case "completion":
		return res.CompletionTokens
	default:
		return res.TotalTokens
	}
}

func buildRecommendationsFromMatches(matches []models.Match, limit int) []TodayRecommendedMatch {
	if limit <= 0 {
		limit = 3
	}
	out := make([]TodayRecommendedMatch, 0, minInt(limit, len(matches)))
	for i, m := range matches {
		if i >= limit {
			break
		}
		rating := m.ImportanceLevel + 2
		if rating > 5 {
			rating = 5
		}
		if rating < 1 {
			rating = 3
		}
		reason := cleanReason(m.RecommendReason)
		if reason == "" {
			reason = "对阵信息明确，适合作为今日观赛候选。"
		}
		out = append(out, TodayRecommendedMatch{
			MatchID:     m.ID,
			Title:       fmt.Sprintf("%s vs %s", safeTeamName(m.HomeTeam), safeTeamName(m.AwayTeam)),
			KickoffTime: m.KickoffTimeUTC.UTC().Format(time.RFC3339),
			Reason:      reason,
			Rating:      rating,
		})
	}
	return out
}

func validateRecommendations(items []TodayRecommendedMatch, matches []models.Match, limit int) []TodayRecommendedMatch {
	if limit <= 0 {
		limit = 3
	}
	matchByID := make(map[uint]models.Match, len(matches))
	for _, m := range matches {
		matchByID[m.ID] = m
	}
	out := make([]TodayRecommendedMatch, 0, minInt(limit, len(items)))
	seen := map[uint]bool{}
	for _, item := range items {
		if len(out) >= limit || item.MatchID == 0 || seen[item.MatchID] {
			continue
		}
		m, ok := matchByID[item.MatchID]
		if !ok {
			continue
		}
		seen[item.MatchID] = true
		if strings.TrimSpace(item.Title) == "" {
			item.Title = fmt.Sprintf("%s vs %s", safeTeamName(m.HomeTeam), safeTeamName(m.AwayTeam))
		}
		if strings.TrimSpace(item.KickoffTime) == "" {
			item.KickoffTime = m.KickoffTimeUTC.UTC().Format(time.RFC3339)
		}
		item.Rating = ai.ClampInt(item.Rating, 1, 5)
		out = append(out, item)
	}
	return out
}

func buildGroupTeams(standings []models.GroupStanding) []GroupAnalysisTeam {
	out := make([]GroupAnalysisTeam, 0, len(standings))
	for _, s := range standings {
		note := fmt.Sprintf("%d 分，净胜球 %d", s.Points, s.GoalDifference)
		out = append(out, GroupAnalysisTeam{
			TeamID:   s.TeamID,
			TeamName: safeTeamName(s.Team),
			Status:   s.QualificationStatus,
			Note:     note,
		})
	}
	return out
}

func conversationResponse(conv models.AIConversation, messages []models.AIMessage) AIConversationResponse {
	out := AIConversationResponse{
		ID:          conv.ID,
		Title:       conv.Title,
		ContextType: conv.ContextType,
		ContextID:   conv.ContextID,
		LastMessage: conv.LastMessage,
		CreatedAt:   conv.CreatedAt,
		UpdatedAt:   conv.UpdatedAt,
	}
	if messages != nil {
		out.Messages = make([]AIChatMessageResponse, 0, len(messages))
		for _, msg := range messages {
			out.Messages = append(out.Messages, messageResponse(msg))
		}
	}
	return out
}

func messageResponse(msg models.AIMessage) AIChatMessageResponse {
	return AIChatMessageResponse{ID: msg.ID, Role: msg.Role, Content: msg.Content, CreatedAt: msg.CreatedAt}
}

func conversationTitle(message string) string {
	runes := []rune(strings.TrimSpace(message))
	if len(runes) > 24 {
		runes = runes[:24]
	}
	title := strings.TrimSpace(string(runes))
	if title == "" {
		return "新会话"
	}
	return title
}

func defaultContext(contextType string) string {
	switch contextType {
	case "match", "team", "group":
		return contextType
	default:
		return "general"
	}
}

func userKey(userID *uint) string {
	if userID == nil || *userID == 0 {
		return "guest"
	}
	return fmt.Sprintf("%d", *userID)
}

func safeTeamName(t models.Team) string {
	if strings.TrimSpace(t.Name) != "" && !strings.Contains(t.Name, "锟") && !strings.Contains(t.Name, "涓") {
		return t.Name
	}
	if strings.TrimSpace(t.NameEn) != "" {
		return t.NameEn
	}
	if strings.TrimSpace(t.FIFACode) != "" {
		return t.FIFACode
	}
	return "TBD"
}

func cleanReason(s string) string {
	if strings.Contains(s, "锟") || strings.Contains(s, "涓") || strings.Contains(s, "鍦") {
		return ""
	}
	return strings.TrimSpace(s)
}

func nonNilStrings(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func aiUserError(err error) error {
	if err == nil {
		return fmt.Errorf("AI 服务暂时不可用")
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "not configured") || strings.Contains(msg, "ai_api_key") {
		return fmt.Errorf("AI 服务未配置，请先设置 AI_API_KEY")
	}
	return fmt.Errorf("AI 服务暂时不可用")
}
