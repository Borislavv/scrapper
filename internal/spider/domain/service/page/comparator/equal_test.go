package pagecomparator

import (
	"regexp"
	"testing"
)

// TestDateRegex checks the regex for date replacing.
func TestDateRegex(t *testing.T) {
	dates := []string{
		"31.12.2023",       // d.m.Y
		"2023.12.31",       // Y.m.d
		"2023-12-31",       // Y-m-d
		"2023/12/31",       // Y/m/d
		"2023.12.31 14:30", // Y.m.d H:i
		"2023-12-31 14:30", // Y-m-d H:i
		"2023/12/31 14:30", // Y/m/d H:i
		"31.12.2023",       // d.m.Y
		"31-12-2023",       // d-m-Y
		"31/12/2023",       // d/m/Y
		"31.12.2023 14:30", // d.m.Y H:i
		"31-12-2023 14:30", // d-m-Y H:i
		"31/12/2023 14:30", // d/m/Y H:i
		"23.12.31",         // y.m.d
		"23-12-31",         // y-m-d
		"23/12/31",         // y/m/d
		"23.12.31 14:30",   // y.m.d H:i
		"23-12-31 14:30",   // y-m-d H:i
		"23/12/31 14:30",   // y/m/d H:i
		"12.2023",          // m.Y
		"12-2023",          // m-Y
		"12/2023",          // m/Y
		"1996",             // Y
		"01",               // m
		"29",               // d
	}

	rgx := regexp.MustCompile(dateRgxPatter)

	for _, date := range dates {
		if !rgx.MatchString(date) {
			t.Errorf("date regex failed, date: %s", date)
		}
	}
}
