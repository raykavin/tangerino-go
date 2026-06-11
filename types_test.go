package tangerino_test

import (
	"encoding/json"
	"testing"
	"time"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestUnixMilliTime_Raw(t *testing.T) {
	ts := tangerino.UnixMilliTime(1776049200000)
	if got := ts.Raw(); got != 1776049200000 {
		t.Errorf("Raw: want 1776049200000, got %d", got)
	}
}

func TestUnixMilliTime_Time(t *testing.T) {
	// 2025-04-11 03:00:00 UTC
	ts := tangerino.UnixMilliTime(1744340400000)
	got := ts.Time()
	if got.Year() != 2025 || got.Month() != time.April || got.Day() != 11 {
		t.Errorf("Time: unexpected date %s", got)
	}
}

func TestUnixMilliTime_Format(t *testing.T) {
	// 2025-01-01 12:00:00 UTC
	ts := tangerino.UnixMilliTime(1735732800000)
	got := ts.Format("02/01/2006")
	if got != "01/01/2025" {
		t.Errorf("Format: want %q, got %q", "01/01/2025", got)
	}
}

func TestUnixMilliTime_JSONRoundtrip(t *testing.T) {
	type payload struct {
		Date tangerino.UnixMilliTime `json:"date"`
	}

	raw := `{"date":1776049200000}`
	var p payload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if p.Date.Raw() != 1776049200000 {
		t.Errorf("Raw after unmarshal: want 1776049200000, got %d", p.Date.Raw())
	}

	enc, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(enc) != raw {
		t.Errorf("Marshal roundtrip: want %s, got %s", raw, enc)
	}
}

func TestDayOffset_Raw(t *testing.T) {
	d := tangerino.DayOffset(39600000)
	if got := d.Raw(); got != 39600000 {
		t.Errorf("Raw: want 39600000, got %d", got)
	}
}

func TestDayOffset_Duration(t *testing.T) {
	d := tangerino.DayOffset(39600000)
	want := 11 * time.Hour
	if got := d.Duration(); got != want {
		t.Errorf("Duration: want %s, got %s", want, got)
	}
}

func TestDayOffset_String(t *testing.T) {
	tests := []struct {
		ms   int64
		want string
	}{
		{0, "00:00"},
		{39600000, "11:00"},  // 11 h
		{54000000, "15:00"},  // 15 h
		{61200000, "17:00"},  // 17 h
		{75600000, "21:00"},  // 21 h
		{86400000, "24:00"},  // midnight next day
		{100800000, "28:00"}, // 12x36 shift past midnight
		{50400000, "14:00"},  // 14 h
		{51300000, "14:15"},  // 14 h 15 min
	}

	for _, tt := range tests {
		d := tangerino.DayOffset(tt.ms)
		if got := d.String(); got != tt.want {
			t.Errorf("DayOffset(%d).String(): want %q, got %q", tt.ms, tt.want, got)
		}
	}
}

func TestDayOffset_JSONRoundtrip(t *testing.T) {
	type payload struct {
		Start tangerino.DayOffset `json:"start"`
	}

	raw := `{"start":39600000}`
	var p payload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if p.Start.Raw() != 39600000 {
		t.Errorf("Raw after unmarshal: want 39600000, got %d", p.Start.Raw())
	}
	if p.Start.String() != "11:00" {
		t.Errorf("String after unmarshal: want %q, got %q", "11:00", p.Start.String())
	}
}
