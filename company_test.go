package tangerino_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestCompaniesService_List(t *testing.T) {
	want := tangerino.Page[tangerino.Company]{
		Content: []tangerino.Company{
			{
				ID:              2183749,
				CNPJ:            "45.959.037/0001-80",
				ExternalID:      "45959037",
				SocialReason:    "CS SERVIÇOS LTDA",
				FantasyName:     "CS SERVIÇOS",
				DescriptionName: "CS SERVIÇOS",
			},
			{
				ID:              2183763,
				CNPJ:            "36.374.509/0001-41",
				SocialReason:    "JR SERVIÇOS LTDA",
				FantasyName:     "JR SERVIÇOS",
				DescriptionName: "JR SERVIÇOS",
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
		if r.URL.Path != "/companies" {
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

	page, err := client.Companies.List(context.Background(), tangerino.ListCompaniesParams{
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

	c0 := page.Content[0]
	if c0.ID != want.Content[0].ID {
		t.Errorf("ID: want %d, got %d", want.Content[0].ID, c0.ID)
	}
	if c0.CNPJ != want.Content[0].CNPJ {
		t.Errorf("CNPJ: want %q, got %q", want.Content[0].CNPJ, c0.CNPJ)
	}
	if c0.SocialReason != want.Content[0].SocialReason {
		t.Errorf("SocialReason: want %q, got %q", want.Content[0].SocialReason, c0.SocialReason)
	}

	// Second company has no ExternalID.
	if page.Content[1].ExternalID != "" {
		t.Errorf("ExternalID: want empty, got %q", page.Content[1].ExternalID)
	}
}

func TestCompaniesService_List_NoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("expected empty query string, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tangerino.Page[tangerino.Company]{})
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.Companies.List(context.Background(), tangerino.ListCompaniesParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestCompaniesService_List_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("bad", "creds")

	_, err := client.Companies.List(context.Background(), tangerino.ListCompaniesParams{})
	if !tangerino.IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got: %v", err)
	}
}
