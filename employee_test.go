package tangerino_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestEmployeesService_List(t *testing.T) {
	birthDate := tangerino.UnixMilliTime(550724400000)

	want := tangerino.Page[tangerino.Employee]{
		Content: []tangerino.Employee{
			{
				ID:            6419626,
				Name:          "CELMA PEREIRA DA SILVA",
				Email:         "Celmasilvademorais125@gmail.com",
				BirthDate:     &birthDate,
				CPF:           "87644894268",
				AdmissionDate: 1776049200000,
				CurrentWorkSchedule: tangerino.WorkScheduleRef{
					ID:        2000000,
					StartDate: 1776049200000,
					Inactive:  false,
				},
				Company:      tangerino.EntityRef{ID: 2189522},
				JobRole:      tangerino.EntityRef{ID: 2292991},
				LastManager:  tangerino.EntityRef{ID: 0},
				Managers:     []tangerino.EntityRef{},
				WorkplaceList: []tangerino.EntityRef{{ID: 2230906}},
				EffectiveDate: 1776049200000,
				Fired:            false,
				CanViewWorkgroup: true,
				Status:           0,
				RecordsPunch:     true,
			},
		},
		First:            true,
		Last:             false,
		TotalElements:    38,
		TotalPages:       2,
		NumberOfElements: 20,
		Size:             20,
		Number:           0,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/employee/find-all" {
			t.Errorf("unexpected path: got %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("pageSize"); got != "20" {
			t.Errorf("unexpected pageSize: got %q", got)
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

	opt, err := tangerino.WithBaseURL(srv.URL)
	if err != nil {
		t.Fatalf("WithBaseURL: %v", err)
	}

	client, err := tangerino.NewClient("user", "pass", opt)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	page, err := client.Employees.List(context.Background(), tangerino.ListEmployeesParams{
		PageSize: 20,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if page.TotalElements != want.TotalElements {
		t.Errorf("TotalElements: want %d, got %d", want.TotalElements, page.TotalElements)
	}
	if page.TotalPages != want.TotalPages {
		t.Errorf("TotalPages: want %d, got %d", want.TotalPages, page.TotalPages)
	}
	if len(page.Content) != len(want.Content) {
		t.Fatalf("Content length: want %d, got %d", len(want.Content), len(page.Content))
	}

	e := page.Content[0]
	if e.ID != want.Content[0].ID {
		t.Errorf("ID: want %d, got %d", want.Content[0].ID, e.ID)
	}
	if e.Name != want.Content[0].Name {
		t.Errorf("Name: want %q, got %q", want.Content[0].Name, e.Name)
	}
	if e.BirthDate == nil || *e.BirthDate != birthDate {
		t.Errorf("BirthDate: want %d, got %v", birthDate, e.BirthDate)
	}
}

func TestEmployeesService_List_NoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("expected empty query string, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tangerino.Page[tangerino.Employee]{}); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	opt, _ := tangerino.WithBaseURL(srv.URL)
	client, _ := tangerino.NewClient("user", "pass", opt)

	_, err := client.Employees.List(context.Background(), tangerino.ListEmployeesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestEmployeesService_List_Pagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page param: want %q, got %q", "1", got)
		}
		if got := r.URL.Query().Get("size"); got != "10" {
			t.Errorf("size param: want %q, got %q", "10", got)
		}

		page := tangerino.Page[tangerino.Employee]{
			Number:    1,
			Size:      10,
			Last:      true,
			TotalPages: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()

	opt, _ := tangerino.WithBaseURL(srv.URL)
	client, _ := tangerino.NewClient("user", "pass", opt)

	page, err := client.Employees.List(context.Background(), tangerino.ListEmployeesParams{
		Page: 1,
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

func TestPage_HasNext(t *testing.T) {
	middle := tangerino.Page[tangerino.Employee]{Last: false, Number: 0, TotalPages: 3}
	if !middle.HasNext() {
		t.Error("expected HasNext() = true for a non-last page")
	}
	if middle.NextPageNumber() != 1 {
		t.Errorf("NextPageNumber: want 1, got %d", middle.NextPageNumber())
	}

	last := tangerino.Page[tangerino.Employee]{Last: true, Number: 2, TotalPages: 3}
	if last.HasNext() {
		t.Error("expected HasNext() = false for the last page")
	}
	if last.NextPageNumber() != -1 {
		t.Errorf("NextPageNumber: want -1, got %d", last.NextPageNumber())
	}
}
