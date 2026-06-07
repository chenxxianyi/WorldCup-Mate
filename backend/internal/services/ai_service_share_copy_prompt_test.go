package services

import (
	"strings"
	"testing"
)

func TestNormalizeShareCopyRequestFallsBackToDefaults(t *testing.T) {
	got := normalizeShareCopyRequest(ShareCopyRequest{
		MatchID:  9,
		Platform: "unknown",
		Tone:     " PROFESSIONAL ",
		Length:   "",
	})

	if got.MatchID != 9 {
		t.Fatalf("unexpected match id: %d", got.MatchID)
	}
	if got.Platform != "general" {
		t.Fatalf("unexpected platform: %s", got.Platform)
	}
	if got.Tone != "professional" {
		t.Fatalf("unexpected tone: %s", got.Tone)
	}
	if got.Length != "short" {
		t.Fatalf("unexpected length: %s", got.Length)
	}
}

func TestBuildShareCopyPromptKeepsOptionInstructionsOnRetry(t *testing.T) {
	req := normalizeShareCopyRequest(ShareCopyRequest{
		MatchID:  9,
		Platform: "xiaohongshu",
		Tone:     "beginner",
		Length:   "long",
	})

	prompt := buildShareCopyPrompt(req, "Home vs Away", true)
	for _, want := range []string{
		"TASK:share_copy",
		"platform=xiaohongshu",
		"tone=beginner",
		"length=long",
		"This is a retry",
		"Home vs Away",
		`{"title":"","content":"","tips":[]}`,
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}
