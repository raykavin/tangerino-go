package tangerino_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tangerino "github.com/raykavin/tangerino-go"
)

// realPunchesJSON is a subset of the actual API response used across tests.
const realPunchesJSON = `[
  {
    "id": 1667647289,
    "date": "2026-06-01",
    "startDate": "2026-06-01T14:06:00",
    "endDate": "2026-06-01T18:02:00",
    "status": 2,
    "startManual": false,
    "endManual": false,
    "pending": false,
    "accredited": false,
    "adjustment": false,
    "canceled": false,
    "totalHours": 14160000
  },
  {
    "id": 1667828834,
    "date": "2026-06-10",
    "startDate": "2026-06-10T18:07:00",
    "status": 2,
    "startManual": false,
    "endManual": false,
    "pending": false,
    "accredited": false,
    "adjustment": false,
    "canceled": false,
    "totalHours": 0
  }
]`

func TestPunchesService_List(t *testing.T) {
	adj := true
	pending := false

	startDate := time.Unix(1772334000, 0)
	endDate := time.Unix(1776199815, 0)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/punch/v2/punches/employees/6419626" {
			t.Errorf("unexpected path: got %q", r.URL.Path)
		}

		q := r.URL.Query()
		if got := q.Get("status"); got != "3" {
			t.Errorf("status param: want %q, got %q", "3", got)
		}
		if got := q.Get("adjustment"); got != "true" {
			t.Errorf("adjustment param: want %q, got %q", "true", got)
		}
		if got := q.Get("pending"); got != "false" {
			t.Errorf("pending param: want %q, got %q", "false", got)
		}
		if got := q.Get("startDate"); got != "1772334000" {
			t.Errorf("startDate param: want %q, got %q", "1772334000", got)
		}
		if got := q.Get("endDate"); got != "1776199815" {
			t.Errorf("endDate param: want %q, got %q", "1776199815", got)
		}
		if _, _, ok := r.BasicAuth(); !ok {
			t.Error("missing Basic Auth header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(realPunchesJSON))
	}))
	defer srv.Close()

	client, err := tangerino.NewClient("user", "pass")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	punches, err := client.Punches.List(context.Background(), 6419626, tangerino.PunchesParams{
		Status:     3,
		Adjustment: &adj,
		StartDate:  startDate,
		EndDate:    endDate,
		Pending:    &pending,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(punches) != 2 {
		t.Fatalf("length: want 2, got %d", len(punches))
	}

	p0 := punches[0]
	if p0.ID != 1667647289 {
		t.Errorf("ID: want 1667647289, got %d", p0.ID)
	}
	if p0.Date != "2026-06-01" {
		t.Errorf("Date: want %q, got %q", "2026-06-01", p0.Date)
	}
	if p0.StartDate.String() != "2026-06-01T14:06:00" {
		t.Errorf("StartDate: want %q, got %q", "2026-06-01T14:06:00", p0.StartDate.String())
	}
	if p0.EndDate == nil {
		t.Fatal("EndDate: want non-nil for closed punch, got nil")
	}
	if p0.EndDate.String() != "2026-06-01T18:02:00" {
		t.Errorf("EndDate: want %q, got %q", "2026-06-01T18:02:00", p0.EndDate.String())
	}
	if p0.Status != 2 {
		t.Errorf("Status: want 2, got %d", p0.Status)
	}
	if p0.TotalHours != 14160000 {
		t.Errorf("TotalHours: want 14160000, got %d", p0.TotalHours)
	}
	if p0.StartManual {
		t.Error("StartManual: want false, got true")
	}

	// Second record has no endDate (open punch).
	p1 := punches[1]
	if p1.EndDate != nil {
		t.Errorf("EndDate: want nil for open punch, got %v", p1.EndDate)
	}
	if p1.TotalHours != 0 {
		t.Errorf("TotalHours: want 0 for open punch, got %d", p1.TotalHours)
	}
}

func TestPunchesService_List_NoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("expected empty query string, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.Punches.List(context.Background(), 1, tangerino.PunchesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestPunchesService_List_TimeConversion(t *testing.T) {
	// Validates that time.Time values are converted to Unix seconds (not millis).
	ts := time.Unix(1700000000, 0)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("startDate"); got != "1700000000" {
			t.Errorf("startDate: want Unix seconds %q, got %q", "1700000000", got)
		}
		if got := r.URL.Query().Get("endDate"); got != "1700000000" {
			t.Errorf("endDate: want Unix seconds %q, got %q", "1700000000", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.Punches.List(context.Background(), 1, tangerino.PunchesParams{
		StartDate: ts,
		EndDate:   ts,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestPunchesService_List_ManualFlags(t *testing.T) {
	body := `[{"id":1,"date":"2026-06-04","startDate":"2026-06-04T08:01:00","endDate":"2026-06-04T12:03:00","status":2,"startManual":true,"endManual":true,"pending":false,"accredited":false,"adjustment":false,"canceled":false,"totalHours":14520000}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	punches, err := client.Punches.List(context.Background(), 1, tangerino.PunchesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(punches) != 1 {
		t.Fatalf("length: want 1, got %d", len(punches))
	}

	p := punches[0]
	if !p.StartManual {
		t.Error("StartManual: want true, got false")
	}
	if !p.EndManual {
		t.Error("EndManual: want true, got false")
	}
	if p.TotalHours != 14520000 {
		t.Errorf("TotalHours: want 14520000, got %d", p.TotalHours)
	}
}

func TestPunchesService_List_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("bad", "creds")

	_, err := client.Punches.List(context.Background(), 1, tangerino.PunchesParams{})
	if !tangerino.IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got: %v", err)
	}
}

func TestLocalDateTime_Parsing(t *testing.T) {
	body := `[{"id":1,"date":"2026-06-01","startDate":"2026-06-01T08:05:00","endDate":"2026-06-01T12:01:00","status":2,"startManual":false,"endManual":false,"pending":false,"accredited":false,"adjustment":false,"canceled":false,"totalHours":14160000}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	punches, err := client.Punches.List(context.Background(), 1, tangerino.PunchesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	p := punches[0]
	start := p.StartDate.Time()

	if start.Year() != 2026 || start.Month() != 6 || start.Day() != 1 {
		t.Errorf("StartDate date part: want 2026-06-01, got %v", start)
	}
	if start.Hour() != 8 || start.Minute() != 5 {
		t.Errorf("StartDate time part: want 08:05, got %02d:%02d", start.Hour(), start.Minute())
	}

	end := p.EndDate.Time()
	if end.Hour() != 12 || end.Minute() != 1 {
		t.Errorf("EndDate time part: want 12:01, got %02d:%02d", end.Hour(), end.Minute())
	}
}
