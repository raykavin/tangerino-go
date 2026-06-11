package tangerino_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestHolidayCalendarsService_List(t *testing.T) {
	want := []tangerino.HolidayCalendar{
		{
			ID:   2543940,
			Name: "CALENDÁRIO PADRÃO 2023",
			Year: 2023,
			Holidays: []tangerino.Holiday{
				{ID: 7628229, Description: "Ano Novo", Date: "2023-01-01"},
				{ID: 7628378, Description: "Carnaval", Date: "2023-02-21"},
			},
		},
		{
			ID:          2669323,
			Name:        "CALENDÁRIO PADRÃO 2024",
			Description: "CALENDÁRIO PADRÃO 2024",
			Year:        2024,
			Holidays: []tangerino.Holiday{
				{ID: 7989779, Description: "Ano Novo", Date: "2024-01-01"},
			},
		},
	}

	envelope := map[string]any{
		"code":     200,
		"status":   "OK",
		"location": "/holiday-calendar/",
		"messages": []string{},
		"item":     want,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/holiday-calendar/" {
			t.Errorf("unexpected path: got %q", r.URL.Path)
		}
		if _, _, ok := r.BasicAuth(); !ok {
			t.Error("missing Basic Auth header")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(envelope); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	client, err := tangerino.NewClient("user", "pass")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	got, err := client.HolidayCalendars.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d calendar(s), got %d", len(want), len(got))
	}

	c0 := got[0]
	if c0.ID != want[0].ID {
		t.Errorf("ID: want %d, got %d", want[0].ID, c0.ID)
	}
	if c0.Year != want[0].Year {
		t.Errorf("Year: want %d, got %d", want[0].Year, c0.Year)
	}
	if len(c0.Holidays) != len(want[0].Holidays) {
		t.Fatalf("Holidays length: want %d, got %d", len(want[0].Holidays), len(c0.Holidays))
	}
	if c0.Holidays[0].Date != want[0].Holidays[0].Date {
		t.Errorf("Holiday.Date: want %q, got %q", want[0].Holidays[0].Date, c0.Holidays[0].Date)
	}

	c1 := got[1]
	if c1.Description != want[1].Description {
		t.Errorf("Description: want %q, got %q", want[1].Description, c1.Description)
	}
}

func TestHolidayCalendarsService_List_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("bad", "creds")

	_, err := client.HolidayCalendars.List(context.Background())
	if !tangerino.IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got: %v", err)
	}
}

func TestHolidayCalendarsService_List_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.HolidayCalendars.List(context.Background())
	if !tangerino.IsServerError(err) {
		t.Errorf("expected server error, got: %v", err)
	}
}
