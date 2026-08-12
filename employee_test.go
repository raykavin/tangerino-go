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
				Company:          tangerino.EntityRef{ID: 2189522},
				JobRole:          tangerino.EntityRef{ID: 2292991},
				LastManager:      tangerino.EntityRef{ID: 0},
				Managers:         []tangerino.EntityRef{},
				WorkplaceList:    []tangerino.EntityRef{{ID: 2230906}},
				EffectiveDate:    1776049200000,
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

	page, err := client.Employees.List(context.Background(), tangerino.ListEmployeesParams{
		Size: 20,
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

	client, _ := tangerino.NewClient("user", "pass")

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
			Number:     1,
			Size:       10,
			Last:       true,
			TotalPages: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

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

const findEmployeeResponse = `{
    "id": 3450546,
    "name": "BISMARK VERAS LEAL",
    "email": "lealbismark89@gmail.com",
    "birthDate": "1989-05-03",
    "phone": "94992398398",
    "cpf": "98386140259",
    "ctps": "94186",
    "series": "00049",
    "pis": "20981163631",
    "admissionDate": "2022-07-22",
    "currentWorkSchedule": {
        "id": 2266697,
        "inactive": false
    },
    "company": {
        "id": 2184050,
        "externalId": "45959167"
    },
    "jobRoleDTO": {
        "id": 3042884,
        "description": "TECNICO DE TELECOMUNICACOES",
        "alterationDate": 1781265924002
    },
    "managers": [
        {
            "id": 2161524,
            "employee": {
                "id": 3550777,
                "name": "ANTONIO WILLIAM MORAES DE OLIVEIRA",
                "admissionDate": "2023-10-13",
                "fired": false,
                "canViewWorkgroup": false,
                "status": 0,
                "doubleBindEmployee": false,
                "recordsPunch": false
            }
        }
    ],
    "workplaceList": [
        {
            "id": 2224274,
            "name": "MARABÁ",
            "active": true
        }
    ],
    "effectiveDate": 1692586800000,
    "externalId": "7",
    "fired": false,
    "state": "Pará",
    "canViewWorkgroup": false,
    "status": 0,
    "doubleBindEmployee": false,
    "recordsPunch": true
}`

func TestEmployeesService_Find(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/employee/find" {
			t.Errorf("unexpected path: got %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("externalId"); got != "string" {
			t.Errorf("externalId param: want %q, got %q", "string", got)
		}
		if got := r.URL.Query().Get("ignoreFired"); got != "true" {
			t.Errorf("ignoreFired param: want %q, got %q", "true", got)
		}
		if got := r.URL.Query().Get("tangerinoId"); got != "7226" {
			t.Errorf("tangerinoId param: want %q, got %q", "7226", got)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(findEmployeeResponse)); err != nil {
			t.Errorf("writing response: %v", err)
		}
	}))
	defer srv.Close()

	client, err := tangerino.NewClient("user", "pass")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	employee, err := client.Employees.Find(context.Background(), tangerino.FindEmployeeParams{
		ExternalID:  "string",
		IgnoreFired: true,
		TangerinoID: 7226,
	})
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	if employee.ID != 3450546 {
		t.Errorf("ID: want %d, got %d", 3450546, employee.ID)
	}
	if employee.Name != "BISMARK VERAS LEAL" {
		t.Errorf("Name: want %q, got %q", "BISMARK VERAS LEAL", employee.Name)
	}
	if employee.BirthDate == nil || *employee.BirthDate != "1989-05-03" {
		t.Errorf("BirthDate: want %q, got %v", "1989-05-03", employee.BirthDate)
	}
	if employee.AdmissionDate != "2022-07-22" {
		t.Errorf("AdmissionDate: want %q, got %q", "2022-07-22", employee.AdmissionDate)
	}
	if employee.Phone != "94992398398" {
		t.Errorf("Phone: want %q, got %q", "94992398398", employee.Phone)
	}
	if employee.CTPS != "94186" {
		t.Errorf("CTPS: want %q, got %q", "94186", employee.CTPS)
	}
	if employee.Series != "00049" {
		t.Errorf("Series: want %q, got %q", "00049", employee.Series)
	}
	if employee.State != "Pará" {
		t.Errorf("State: want %q, got %q", "Pará", employee.State)
	}
	if employee.Company.ID != 2184050 || employee.Company.ExternalID != "45959167" {
		t.Errorf("Company: want {2184050 45959167}, got %+v", employee.Company)
	}
	if employee.JobRole.Description != "TECNICO DE TELECOMUNICACOES" {
		t.Errorf("JobRole.Description: want %q, got %q", "TECNICO DE TELECOMUNICACOES", employee.JobRole.Description)
	}
	if employee.JobRole.AlterationDate != 1781265924002 {
		t.Errorf("JobRole.AlterationDate: want %d, got %d", 1781265924002, employee.JobRole.AlterationDate)
	}
	if len(employee.Managers) != 1 || employee.Managers[0].Employee.Name != "ANTONIO WILLIAM MORAES DE OLIVEIRA" {
		t.Errorf("Managers: unexpected value %+v", employee.Managers)
	}
	if len(employee.WorkplaceList) != 1 || employee.WorkplaceList[0].Name != "MARABÁ" || !employee.WorkplaceList[0].Active {
		t.Errorf("WorkplaceList: unexpected value %+v", employee.WorkplaceList)
	}
	if employee.EffectiveDate != 1692586800000 {
		t.Errorf("EffectiveDate: want %d, got %d", 1692586800000, employee.EffectiveDate)
	}
	if employee.ExternalID != "7" {
		t.Errorf("ExternalID: want %q, got %q", "7", employee.ExternalID)
	}
	if !employee.RecordsPunch {
		t.Error("RecordsPunch: want true, got false")
	}
}

func TestEmployeesService_Find_NoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("expected empty query string, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tangerino.EmployeeDetail{}); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	client, _ := tangerino.NewClient("user", "pass")

	_, err := client.Employees.Find(context.Background(), tangerino.FindEmployeeParams{})
	if err != nil {
		t.Fatalf("Find: %v", err)
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
