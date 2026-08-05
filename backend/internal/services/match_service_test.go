package services

import (
	"testing"
	"time"
)

// DATA-05E: calendar-day boundaries in user timezones, including the
// DST-transition day and cross-midnight cases.
func TestDayRangeInTZ(t *testing.T) {
	cases := []struct {
		name      string
		nowUTC    string
		tz        string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "UTC default",
			nowUTC:    "2026-03-08T16:00:00Z",
			tz:        "",
			wantStart: "2026-03-08T00:00:00Z",
			wantEnd:   "2026-03-09T00:00:00Z",
		},
		{
			name:      "Beijing cross-midnight: UTC 16:00 is next calendar day",
			nowUTC:    "2026-03-08T16:00:00Z", // 2026-03-09 00:00 +08
			tz:        "Asia/Shanghai",
			wantStart: "2026-03-08T16:00:00Z",
			wantEnd:   "2026-03-09T16:00:00Z",
		},
		{
			name:      "New York winter day (UTC-5)",
			nowUTC:    "2026-03-01T12:00:00Z", // 07:00 EST
			tz:        "America/New_York",
			wantStart: "2026-03-01T05:00:00Z",
			wantEnd:   "2026-03-02T05:00:00Z",
		},
		{
			name:      "US DST spring-forward day is 23 hours long",
			nowUTC:    "2026-03-08T05:00:00Z", // 00:00 EST, switch at 07:00Z
			tz:        "America/New_York",
			wantStart: "2026-03-08T05:00:00Z",
			wantEnd:   "2026-03-09T04:00:00Z", // 00:00 EDT the next day
		},
		{
			name:      "invalid tz falls back to UTC",
			nowUTC:    "2026-03-08T16:00:00Z",
			tz:        "Not/AZone",
			wantStart: "2026-03-08T00:00:00Z",
			wantEnd:   "2026-03-09T00:00:00Z",
		},
	}
	for _, c := range cases {
		now, err := time.Parse(time.RFC3339, c.nowUTC)
		if err != nil {
			t.Fatalf("%s: bad fixture %v", c.name, err)
		}
		start, end := dayRangeInTZ(now, c.tz)
		wantStart, _ := time.Parse(time.RFC3339, c.wantStart)
		wantEnd, _ := time.Parse(time.RFC3339, c.wantEnd)
		if !start.Equal(wantStart) || !end.Equal(wantEnd) {
			t.Errorf("%s: got [%s, %s), want [%s, %s)",
				c.name, start.Format(time.RFC3339), end.Format(time.RFC3339),
				wantStart.Format(time.RFC3339), wantEnd.Format(time.RFC3339))
		}
	}
}
