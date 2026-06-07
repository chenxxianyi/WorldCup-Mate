package services

import "testing"

func TestDecodeShareCopyAcceptsNestedResult(t *testing.T) {
	raw := `{"result":{"headline":"今晚看点","text":"这场值得约朋友一起看。","notes":["赛前确认开球时间"]}}`

	got, jsonText, err := decodeShareCopy(raw)
	if err != nil {
		t.Fatalf("decodeShareCopy returned error: %v", err)
	}
	if got.Title != "今晚看点" || got.Content != "这场值得约朋友一起看。" {
		t.Fatalf("unexpected share copy: %+v", got)
	}
	if jsonText == "" {
		t.Fatal("expected normalized json text")
	}
}

func TestDecodeShareCopyAcceptsPlainText(t *testing.T) {
	raw := "今晚这场比赛适合边吃宵夜边看，节奏和对抗都值得期待。"

	got, _, err := decodeShareCopy(raw)
	if err != nil {
		t.Fatalf("decodeShareCopy returned error: %v", err)
	}
	if got.Content != raw {
		t.Fatalf("unexpected content: %q", got.Content)
	}
	if got.Tips == nil {
		t.Fatal("expected tips to be a non-nil slice")
	}
}
