package ai

import (
	"strings"
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
	if err := ValidateOutput("这场比赛不要参考赔率和盘口"); err == nil {
		t.Fatal("expected Chinese betting terms to be rejected")
	}
	if err := ValidateOutput("watch the opening tempo and substitutions"); err != nil {
		t.Fatalf("expected safe output, got %v", err)
	}
}

func TestPromptRequiresContextForLiveFacts(t *testing.T) {
	prompt := BuildSystemPrompt()
	for _, want := range []string{"不能把模型记忆当作最新事实", "不要声称已经联网查询"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt should contain %q, got: %s", want, prompt)
		}
	}
}
