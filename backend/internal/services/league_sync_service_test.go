package services

import "testing"

func TestZoneForPosition(t *testing.T) {
	cases := []struct {
		code string
		pos  int
		want string
	}{
		{"PL", 1, "champions_league"},
		{"PL", 4, "champions_league"},
		{"PL", 5, "europa_league"},
		{"PL", 6, "europa_league"},
		{"PL", 18, "relegation"},
		{"PL", 20, "relegation"},
		{"PL", 10, ""},
		{"PD", 2, "champions_league"},
		{"SA", 19, "relegation"},
		{"BL1", 15, ""},
		{"BL1", 16, "relegation"},
		{"FL1", 3, "champions_league"},
		{"FL1", 4, "europa_league"},
		{"fl1", 4, "europa_league"}, // case-insensitive
		{"WC", 1, ""},               // cup competitions have no zone rules
	}
	for _, c := range cases {
		if got := zoneForPosition(c.code, c.pos); got != c.want {
			t.Errorf("zoneForPosition(%q, %d) = %q, want %q", c.code, c.pos, got, c.want)
		}
	}
}

func TestConfigureLeagueSyncParsing(t *testing.T) {
	ConfigureLeagueSync("PL:2025, pd:2026 ,BL1", 30, "", "key")
	cfg := GetLeagueSyncConfig()
	if len(cfg.Targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(cfg.Targets))
	}
	if cfg.Targets[0].Code != "PL" || cfg.Targets[0].Season != 2025 {
		t.Errorf("target[0] = %+v, want {PL 2025}", cfg.Targets[0])
	}
	if cfg.Targets[1].Code != "PD" || cfg.Targets[1].Season != 2026 {
		t.Errorf("target[1] = %+v, want {PD 2026} (whitespace + lowercase)", cfg.Targets[1])
	}
	if cfg.Targets[2].Code != "BL1" || cfg.Targets[2].Season != 0 {
		t.Errorf("target[2] = %+v, want {BL1 0} (season optional)", cfg.Targets[2])
	}
	if !IsLeagueSyncEnabled() {
		t.Error("expected league sync enabled with targets + api key")
	}

	// Empty config must disable the league sync (legacy behavior).
	ConfigureLeagueSync("", 0, "", "")
	if IsLeagueSyncEnabled() {
		t.Error("empty targets must disable league sync")
	}
}
