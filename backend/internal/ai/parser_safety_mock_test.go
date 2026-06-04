package ai

import (
	"context"
	"testing"
)

func TestDecodeJSONExtractsObject(t *testing.T) {
	type payload struct {
		Summary string   `json:"summary"`
		Items   []string `json:"items"`
	}

	got, _, err := DecodeJSON("```json\n{\"summary\":\"ok\",\"items\":[\"a\"]}\n```", payload{})
	if err != nil {
		t.Fatalf("DecodeJSON returned error: %v", err)
	}
	if got.Summary != "ok" || len(got.Items) != 1 || got.Items[0] != "a" {
		t.Fatalf("unexpected decoded payload: %+v", got)
	}
}

func TestValidateOutputBlocksUnsafeTerms(t *testing.T) {
	if err := ValidateOutput("this mentions odds and betting lines"); err == nil {
		t.Fatal("expected unsafe output to be rejected")
	}
	if err := ValidateOutput("watch the opening tempo and substitutions"); err != nil {
		t.Fatalf("expected safe output, got %v", err)
	}
}

func TestMockProviderMatchInsightJSON(t *testing.T) {
	provider := NewMockProvider("")
	res, err := provider.Chat(context.Background(), ChatRequest{UserPrompt: "TASK:match_insight"})
	if err != nil {
		t.Fatalf("mock provider failed: %v", err)
	}
	type matchInsight struct {
		Summary     string   `json:"summary"`
		WatchRating int      `json:"watch_rating"`
		Reasons     []string `json:"reasons"`
	}
	got, _, err := DecodeJSON(res.Content, matchInsight{})
	if err != nil {
		t.Fatalf("mock response is not valid JSON: %v", err)
	}
	if got.Summary == "" || got.WatchRating == 0 || len(got.Reasons) == 0 {
		t.Fatalf("mock response missing expected fields: %+v", got)
	}
}
