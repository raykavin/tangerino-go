package tangerino

import (
	"context"
	"net/url"
	"strconv"
)

// Page is a generic paginated response from the Tangerino API.
// It wraps a slice of T alongside Spring-style pagination metadata.
// Any endpoint that returns a paginated collection uses this type as its response.
//
// Example:
//
//	page, err := client.Employees.List(ctx, tangerino.ListEmployeesParams{PageSize: 20})
//	for !page.IsLast() {
//	    // process page.Content ...
//	    params.Page++
//	    page, err = client.Employees.List(ctx, params)
//	}
type Page[T any] struct {
	// Content holds the items on the current page.
	Content []T `json:"content"`
	// First indicates whether this is the first page.
	First bool `json:"first"`
	// Last indicates whether this is the last page.
	Last bool `json:"last"`
	// TotalElements is the total number of items across all pages.
	TotalElements int `json:"totalElements"`
	// TotalPages is the total number of pages available.
	TotalPages int `json:"totalPages"`
	// NumberOfElements is the number of items on this page.
	NumberOfElements int `json:"numberOfElements"`
	// Size is the maximum number of items per page as requested.
	Size int `json:"size"`
	// Number is the zero-based index of the current page.
	Number int `json:"number"`
}

// HasNext reports whether there is at least one more page after this one.
func (p *Page[T]) HasNext() bool {
	return !p.Last
}

// NextPageNumber returns the number to pass as Page in the next request.
// Returns -1 when the current page is the last one.
func (p *Page[T]) NextPageNumber() int {
	if p.Last {
		return -1
	}
	return p.Number + 1
}

// EmployeesService handles communication with the employee endpoints.
type EmployeesService struct {
	client *Client
}

// WorkScheduleRef is a lightweight reference to a work schedule as embedded in an employee record.
// For the full work schedule details use WorkSchedule, returned by WorkSchedulesService.
type WorkScheduleRef struct {
	// ID is the unique identifier of the work schedule.
	ID int `json:"id"`
	// StartDate is the Unix timestamp (milliseconds) when the schedule became effective for the employee.
	StartDate UnixMilliTime `json:"startDate"`
	// Inactive indicates whether the schedule has been deactivated.
	Inactive bool `json:"inactive"`
}

// EntityRef is a lightweight reference to a related entity identified by its ID.
// It is used for nested objects where the API returns only the identifier.
type EntityRef struct {
	// ID is the unique identifier of the referenced entity.
	ID int `json:"id"`
}

// Employee represents a single employee record returned by the API.
type Employee struct {
	// ID is the unique identifier of the employee.
	ID int `json:"id"`
	// Name is the employee's full legal name.
	Name string `json:"name"`
	// SocialName is the employee's preferred or social name, if provided.
	SocialName string `json:"socialName"`
	// Email is the employee's contact email address.
	Email string `json:"email"`
	// CPF is the employee's Brazilian tax identification number.
	CPF string `json:"cpf"`
	// PIS is the employee's Social Integration Program number, if provided.
	PIS string `json:"pis"`
	// Gender is the employee's gender as reported by the API (e.g. "MASCULINO", "FEMININO").
	Gender string `json:"gender"`
	// BirthDate is the employee's date of birth as a Unix millisecond timestamp.
	// It is nil when the value is not present in the API response.
	BirthDate *UnixMilliTime `json:"birthDate"`
	// AdmissionDate is the employee's hiring date as a Unix millisecond timestamp.
	AdmissionDate UnixMilliTime `json:"admissionDate"`
	// EffectiveDate is the date the current record became effective, as a Unix millisecond timestamp.
	EffectiveDate UnixMilliTime `json:"effectiveDate"`
	// ExternalID is an optional identifier assigned by an external system.
	ExternalID string `json:"externalId"`
	// CurrentWorkSchedule is a reference to the work schedule currently active for the employee.
	CurrentWorkSchedule WorkScheduleRef `json:"currentWorkSchedule"`
	// Company is a reference to the company the employee belongs to.
	Company EntityRef `json:"company"`
	// JobRole is a reference to the employee's current job role.
	JobRole EntityRef `json:"jobRoleDTO"`
	// LastManager is a reference to the employee's most recent manager.
	LastManager EntityRef `json:"lastManager"`
	// Managers holds references to all current managers for the employee.
	Managers []EntityRef `json:"managers"`
	// WorkplaceList holds references to all workplaces assigned to the employee.
	WorkplaceList []EntityRef `json:"workplaceList"`
	// Fired indicates whether the employee has been terminated.
	Fired bool `json:"fired"`
	// CanViewWorkgroup indicates whether the employee has workgroup visibility permissions.
	CanViewWorkgroup bool `json:"canViewWorkgroup"`
	// Status is the numeric status code for the employee record.
	Status int `json:"status"`
	// DoubleBindEmployee indicates whether the employee is shared across multiple companies.
	DoubleBindEmployee bool `json:"doubleBindEmployee"`
	// RecordsPunch indicates whether the employee uses the punch clock system.
	RecordsPunch bool `json:"recordsPunch"`
}

// ListEmployeesParams holds optional filter and pagination parameters for the employee list endpoint.
// All fields are optional; zero values are omitted from the request.
type ListEmployeesParams struct {
	// BranchExternalID filters employees by the external identifier of their branch.
	BranchExternalID string
	// ManagerExternalID filters employees by the external identifier of their manager.
	ManagerExternalID string
	// LastUpdate filters employees modified after this Unix timestamp in milliseconds.
	LastUpdate int64
	// Page is the zero-based page index to retrieve.
	Page int
	// PageNumber is an alias for Page accepted by the API.
	PageNumber int
	// PageSize is the number of items per page (used alongside PageNumber).
	PageSize int
	// Size is the number of items per page (used alongside Page).
	Size int
	// Offset is the item offset within the result set.
	Offset int
	// ShowFired controls whether terminated employees are included (0 = exclude, 1 = include).
	ShowFired int
}

// List retrieves a single page of employees matching the given parameters.
// All parameters are optional; omit them by using a zero-value ListEmployeesParams.
//
// GET /employee/find-all
func (s *EmployeesService) List(ctx context.Context, params ListEmployeesParams) (*Page[Employee], error) {
	q := url.Values{}

	if params.BranchExternalID != "" {
		q.Set("branchExternalId", params.BranchExternalID)
	}
	if params.ManagerExternalID != "" {
		q.Set("managerExternalId", params.ManagerExternalID)
	}
	if params.LastUpdate != 0 {
		q.Set("lastUpdate", strconv.FormatInt(params.LastUpdate, 10))
	}
	if params.Page != 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.PageNumber != 0 {
		q.Set("pageNumber", strconv.Itoa(params.PageNumber))
	}
	if params.PageSize != 0 {
		q.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.Size != 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}
	if params.Offset != 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.ShowFired != 0 {
		q.Set("showFired", strconv.Itoa(params.ShowFired))
	}

	rawURL := s.client.resolveURL("/employee/find-all")
	if len(q) > 0 {
		rawURL += "?" + q.Encode()
	}

	var page Page[Employee]
	if err := s.client.get(ctx, rawURL, &page); err != nil {
		return nil, err
	}

	return &page, nil
}
