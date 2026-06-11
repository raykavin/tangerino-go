package tangerino_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestWorkplacesService_List(t *testing.T) {
	want := tangerino.Page[tangerino.Workplace]{
		Content: []tangerino.Workplace{
			{
				ID:         2230906,
				Name:       "Sede Central",
				ExternalID: "sede-01",
				Company:    tangerino.EntityRef{ID: 2189522},
				Address:    "Av. Paulista, 1000",
				City:       "São Paulo",
				State:      "SP",
				ZipCode:    "01310-100",
				Inactive:   false,
			},
			{
				ID:       2230907,
				Name:     "Filial Norte",
				Company:  tangerino.EntityRef{ID: 2189522},
				Inactive: false,
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
		if r.URL.Path != "/workplace/find-all" {
			t.Errorf("unexpected path: got %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("size"); got != "20" {
			t.Errorf("unexpected size: got %q", got)
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

	page, err := client.Workplaces.List(context.Background(), tangerino.ListWorkplacesParams{
		Size: 20,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if page.TotalElements != want.TotalElements {
		t.Errorf("TotalElements: want %d, got %d", want.TotalElements, page.TotalElements)
	}
	if len(page.Content) != len(want.Content) {
		t.Fatalf("Content length: want %d, got %d", len(want.Content), len(page.Content))
	}

	w0 := page.Content[0]
	if w0.ID != want.Content[0].ID {
		t.Errorf("ID: want %d, got %d", want.Content[0].ID, w0.ID)
	}
	if w0.Name != want.Content[0].Name {
		t.Errorf("Name: want %q, got %q", want.Content[0].Name, w0.Name)
	}
	if w0.ExternalID != want.Content[0].ExternalID {
		t.Errorf("ExternalID: want %q, got %q", want.Content[0].ExternalID, w0.ExternalID)
	}
	if w0.City != want.Content[0].City {
		t.Errorf("City: want %q, got %q", want.Content[0].City, w0.City)
	}

	// Second workplace has no ExternalID or address fields.
	if page.Content[1].ExternalID != "" {
		t.Errorf("ExternalID: want empty, got %q", page.Content[1].ExternalID)
	}
	if page.Content[1].Address != "" {
		t.Errorf("Address: want empty, got %q", page.Content[1].Address)
	}
}

func TestWorkplacesService_List_NoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("expected empty query string, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tangerino.Page[tangerino.Workplace]{})
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.Workplaces.List(context.Background(), tangerino.ListWorkplacesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestWorkplacesService_List_Pagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page param: want %q, got %q", "1", got)
		}
		if got := r.URL.Query().Get("size"); got != "1" {
			t.Errorf("size param: want %q, got %q", "1", got)
		}

		page := tangerino.Page[tangerino.Workplace]{
			Number:     1,
			Size:       1,
			Last:       true,
			TotalPages: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	// Mirrors the curl example: ?size=1&page=1
	page, err := client.Workplaces.List(context.Background(), tangerino.ListWorkplacesParams{
		Page: 1,
		Size: 1,
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

func TestWorkplacesService_List_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("bad", "creds")

	_, err := client.Workplaces.List(context.Background(), tangerino.ListWorkplacesParams{})
	if !tangerino.IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got: %v", err)
	}
}
