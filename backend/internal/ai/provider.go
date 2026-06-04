package ai

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	SystemPrompt string
	UserPrompt   string
	Messages     []Message
	Model        string
	Temperature  float64
	MaxTokens    int
}

type ChatResponse struct {
	Content          string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type Provider interface {
	Name() string
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

type ProviderConfig struct {
	Provider       string
	BaseURL        string
	APIKey         string
	Model          string
	TimeoutSeconds int
}

func NewProvider(cfg ProviderConfig) (Provider, error) {
	name := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if name == "" || name == "mock" {
		return NewMockProvider(cfg.Model), nil
	}
	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 60
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("AI_API_KEY is required for provider %s", name)
	}
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Name:    name,
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}), nil
}
