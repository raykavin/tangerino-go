package tangerino_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestWorkSchedulesService_List(t *testing.T) {
	endMain := tangerino.DayOffset(61200000)
	endShift1 := tangerino.DayOffset(54000000)
	startShift2 := tangerino.DayOffset(61200000)
	endShift2 := tangerino.DayOffset(75600000)

	want := tangerino.Page[tangerino.WorkSchedule]{
		Content: []tangerino.WorkSchedule{
			{
				ID:       2000000,
				Name:     "Escala Padrão - Tangerino",
				Standard: true,
				Timetable: []tangerino.WorkScheduleTimetable{
					{
						ID:                       2000000,
						Day:                      2,
						StartMainInterval:        54000000,
						EndMainInterval:          &endMain,
						StartShift1:              39600000,
						EndShift1:                &endShift1,
						StartShift2:              &startShift2,
						EndShift2:                &endShift2,
						IntervalPreAssigned1And2: false,
						IntervalPreAssigned2And3: false,
					},
				},
				AlterationDate:          1469744098490,
				PreAssignedInterval:     false,
				ShowIntradayInTimeSheet: false,
				IgnoreHoliday:           false,
				Inactive:                false,
			},
			{
				ID:        2000001,
				Name:      "Escala Padrão - Tangerino Intermitente",
				Standard:  false,
				Timetable: []tangerino.WorkScheduleTimetable{},
			},
		},
		First:            true,
		Last:             true,
		TotalElements:    2,
		TotalPages:       1,
		NumberOfElements: 2,
		Size:             20,
		Number:           0,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/work-schedule" {
			t.Errorf("unexpected path: got %q", r.URL.Path)
		}
		if _, _, ok := r.BasicAuth(); !ok {
			t.Error("missing Basic Auth header")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	client, err := tangerino.NewClient("user", "pass")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	page, err := client.WorkSchedules.List(context.Background(), tangerino.ListWorkSchedulesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(page.Content) != len(want.Content) {
		t.Fatalf("Content length: want %d, got %d", len(want.Content), len(page.Content))
	}

	ws := page.Content[0]
	if ws.ID != want.Content[0].ID {
		t.Errorf("ID: want %d, got %d", want.Content[0].ID, ws.ID)
	}
	if ws.Name != want.Content[0].Name {
		t.Errorf("Name: want %q, got %q", want.Content[0].Name, ws.Name)
	}
	if !ws.Standard {
		t.Error("Standard: want true, got false")
	}
	if len(ws.Timetable) != 1 {
		t.Fatalf("Timetable length: want 1, got %d", len(ws.Timetable))
	}

	tt := ws.Timetable[0]
	if tt.Day != 2 {
		t.Errorf("Day: want 2, got %d", tt.Day)
	}
	if tt.EndMainInterval == nil || *tt.EndMainInterval != endMain {
		t.Errorf("EndMainInterval: want %d, got %v", endMain, tt.EndMainInterval)
	}
	if tt.StartShift2 == nil || *tt.StartShift2 != startShift2 {
		t.Errorf("StartShift2: want %d, got %v", startShift2, tt.StartShift2)
	}

	// Second schedule has empty timetable.
	if len(page.Content[1].Timetable) != 0 {
		t.Errorf("second schedule Timetable: want empty, got %d entries", len(page.Content[1].Timetable))
	}
}

func TestWorkSchedulesService_List_OptionalFields(t *testing.T) {
	// Timetable entry with only mandatory fields (no end times, no second shift).
	body := `{"content":[{"id":1,"name":"Escala","standard":false,"workScheduleTimetableList":[{"id":10,"day":1,"startMainInterval":57600000,"startShift1":43200000,"endShift1":57600000,"intervalPreAssigned1And2":false,"intervalPreAssigned2And3":false,"intervalPreAssigned3And4":false,"intervalPreAssigned4And5":false,"intervalPreAssigned5And6":false}],"alterationDate":1000000,"preAssignedInterval":false,"showIntradayInTimeSheet":false,"ignoreHoliday":false,"inactive":false}],"first":true,"last":true,"totalElements":1,"totalPages":1,"numberOfElements":1,"size":20,"number":0}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	page, err := client.WorkSchedules.List(context.Background(), tangerino.ListWorkSchedulesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	tt := page.Content[0].Timetable[0]
	if tt.EndMainInterval != nil {
		t.Errorf("EndMainInterval: want nil for absent field, got %v", tt.EndMainInterval)
	}
	if tt.StartShift2 != nil {
		t.Errorf("StartShift2: want nil for absent field, got %v", tt.StartShift2)
	}
	if tt.EndShift2 != nil {
		t.Errorf("EndShift2: want nil for absent field, got %v", tt.EndShift2)
	}
}

func TestWorkSchedulesService_List_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("bad", "creds")

	_, err := client.WorkSchedules.List(context.Background(), tangerino.ListWorkSchedulesParams{})
	if !tangerino.IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got: %v", err)
	}
}

func TestWorkSchedulesService_List_NoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("expected empty query string, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tangerino.Page[tangerino.WorkSchedule]{})
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.WorkSchedules.List(context.Background(), tangerino.ListWorkSchedulesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestWorkSchedulesService_List_Pagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page param: want %q, got %q", "2", got)
		}
		if got := r.URL.Query().Get("size"); got != "10" {
			t.Errorf("size param: want %q, got %q", "10", got)
		}

		page := tangerino.Page[tangerino.WorkSchedule]{
			Number:     2,
			Size:       10,
			Last:       true,
			TotalPages: 3,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	page, err := client.WorkSchedules.List(context.Background(), tangerino.ListWorkSchedulesParams{
		Page: 2,
		Size: 10,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if page.HasNext() {
		t.Error("expected HasNext() = false for last page")
	}
	if page.NextPageNumber() != -1 {
		t.Errorf("NextPageNumber: want -1, got %d", page.NextPageNumber())
	}
}
