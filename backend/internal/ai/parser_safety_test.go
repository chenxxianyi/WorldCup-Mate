package ai

import "testing"

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
