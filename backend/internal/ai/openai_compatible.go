package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAICompatibleConfig struct {
	Name    string
	BaseURL string
	APIKey  string
	Model   string
	Timeout time.Duration
}

type OpenAICompatibleProvider struct {
	name   string
	base   string
	key    string
	model  string
	client *http.Client
}

func NewOpenAICompatibleProvider(cfg OpenAICompatibleConfig) *OpenAICompatibleProvider {
	name := strings.TrimSpace(cfg.Name)
	if name == "" {
		name = "openai-compatible"
	}
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = "gpt-4o-mini"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &OpenAICompatibleProvider{
		name:   name,
		base:   base,
		key:    cfg.APIKey,
		model:  model,
		client: &http.Client{Timeout: timeout},
	}
}

func (p *OpenAICompatibleProvider) Name() string {
	return p.name
}

func (p *OpenAICompatibleProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if strings.TrimSpace(p.key) == "" {
		return nil, fmt.Errorf("AI provider is not configured")
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = p.model
	}

	messages := make([]Message, 0, len(req.Messages)+2)
	if strings.TrimSpace(req.SystemPrompt) != "" {
		messages = append(messages, Message{Role: "system", Content: req.SystemPrompt})
	}
	messages = append(messages, req.Messages...)
	if strings.TrimSpace(req.UserPrompt) != "" {
		messages = append(messages, Message{Role: "user", Content: req.UserPrompt})
	}

	payload := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode AI request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create AI request")
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.key)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("AI provider request failed")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("AI provider returned status %d", resp.StatusCode)
	}

	var parsed struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse AI response")
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("AI provider returned empty response")
	}
	if parsed.Model == "" {
		parsed.Model = model
	}
	return &ChatResponse{
		Content:          parsed.Choices[0].Message.Content,
		Model:            parsed.Model,
		PromptTokens:     parsed.Usage.PromptTokens,
		CompletionTokens: parsed.Usage.CompletionTokens,
		TotalTokens:      parsed.Usage.TotalTokens,
	}, nil
}

func (p *OpenAICompatibleProvider) ChatStream(ctx context.Context, req ChatRequest, cb StreamCallback) (*ChatResponse, error) {
	if strings.TrimSpace(p.key) == "" {
		return nil, fmt.Errorf("AI provider is not configured")
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = p.model
	}

	messages := make([]Message, 0, len(req.Messages)+2)
	if strings.TrimSpace(req.SystemPrompt) != "" {
		messages = append(messages, Message{Role: "system", Content: req.SystemPrompt})
	}
	messages = append(messages, req.Messages...)
	if strings.TrimSpace(req.UserPrompt) != "" {
		messages = append(messages, Message{Role: "user", Content: req.UserPrompt})
	}

	payload := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":      true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode AI request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create AI request")
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.key)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("AI provider request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 2<<20))
		return nil, fmt.Errorf("AI provider returned status %d", resp.StatusCode)
	}

	var full strings.Builder
	final := &ChatResponse{Model: model}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			final.Content = full.String()
			if cb != nil {
				if err := cb(StreamDelta{Done: true, Response: final}); err != nil {
					return final, err
				}
			}
			return final, nil
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
			Usage struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			} `json:"usage"`
			Model string `json:"model"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if chunk.Model != "" {
			final.Model = chunk.Model
		}
		if chunk.Usage.TotalTokens > 0 || chunk.Usage.PromptTokens > 0 || chunk.Usage.CompletionTokens > 0 {
			final.PromptTokens = chunk.Usage.PromptTokens
			final.CompletionTokens = chunk.Usage.CompletionTokens
			final.TotalTokens = chunk.Usage.TotalTokens
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}
		full.WriteString(content)
		if cb != nil {
			if err := cb(StreamDelta{Content: content, Response: final}); err != nil {
				return final, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return final, fmt.Errorf("AI provider stream failed")
	}
	final.Content = full.String()
	if cb != nil {
		if err := cb(StreamDelta{Done: true, Response: final}); err != nil {
			return final, err
		}
	}
	return final, nil
}
