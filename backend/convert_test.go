package backend

import (
	"strings"
	"testing"
)

func TestConvert_Units(t *testing.T) {
	tests := []struct {
		query       string
		wantSuccess bool
		wantCat     string
		checkSubstr string
	}{
		{"10 km to mi", true, "length", "mi"},
		{"100 m in ft", true, "length", "ft"},
		{"50 kg to lb", true, "mass", "lb"},
		{"16 gb in mb", true, "data", "16384"},
		{"100 km/h to mph", true, "speed", "mph"},
		{"2 hours to minutes", true, "time", "120"},
		{"100 c to f", true, "temperature", "212"},
		{"32 f to c", true, "temperature", "0.00"},
		{"16px to rem", true, "css", "1rem"},
		{"2rem to px", true, "css", "32px"},
		{"100 usd to eur", true, "currency", "EUR"},
		{"#3875f6", true, "color", "rgb(56, 117, 246)"},
		{"rgb(56, 117, 246)", true, "color", "#3875F6"},
		{"now to epoch", true, "time", "epoch"},
		{"epoch 1700000000", true, "time", "UTC"},
		{"1700000000 to date", true, "time", "UTC"},
		{"days in march", true, "calendar", "31 days"},
		{"days in current month", true, "calendar", "days"},
		{"days in february 2024", true, "calendar", "29 days"},
		{"days in feb 2025", true, "calendar", "28 days"},
		{"days in this month", true, "calendar", "days"},
		{"what is the weekday?", true, "calendar", "Today is"},
		{"what day is today", true, "calendar", "Today is"},
		{"weekday", true, "calendar", "Today is"},
		{"weekday of 2026-10-14", true, "calendar", "Wednesday"},
		{"today", true, "calendar", "20"},
		{"week number", true, "calendar", "Week"},
		{"day of year", true, "calendar", "Day"},
		{"days in 2024", true, "calendar", "366 days"},
		{"days in 2025", true, "calendar", "365 days"},
		{"invalid query", false, "", ""},
		{"", false, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			res := Convert(tc.query)
			if res.Success != tc.wantSuccess {
				t.Fatalf("for %q, expected success=%v, got %v", tc.query, tc.wantSuccess, res.Success)
			}
			if tc.wantSuccess {
				if tc.wantCat != "" && res.Category != tc.wantCat {
					t.Errorf("expected category %q, got %q", tc.wantCat, res.Category)
				}
				if tc.checkSubstr != "" && !strings.Contains(res.Formatted, tc.checkSubstr) && !strings.Contains(res.Result, tc.checkSubstr) {
					t.Errorf("expected result to contain %q, got Formatted=%q Result=%q", tc.checkSubstr, res.Formatted, res.Result)
				}
			}
		})
	}
}
