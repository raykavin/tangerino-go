package tangerino

import (
	"context"
	"net/url"
	"strconv"
)

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
	// Size is the number of items per page.
	Size int
	// Offset is the item offset within the result set.
	Offset int
	// ShowFired controls whether terminated employees are included (0 = exclude, 1 = include).
	ShowFired int
}

// CompanyRef is a lightweight reference to the company an employee belongs to,
// as returned by the single-employee lookup endpoint.
type CompanyRef struct {
	// ID is the unique identifier of the company.
	ID int `json:"id"`
	// ExternalID is an optional identifier assigned by an external system.
	ExternalID string `json:"externalId"`
}

// JobRoleRef is a reference to an employee's job role, as returned by the
// single-employee lookup endpoint.
type JobRoleRef struct {
	// ID is the unique identifier of the job role.
	ID int `json:"id"`
	// Description is the human-readable name of the job role.
	Description string `json:"description"`
	// AlterationDate is when the job role was last changed, as a Unix millisecond timestamp.
	AlterationDate UnixMilliTime `json:"alterationDate"`
}

// WorkplaceRef is a reference to a workplace assigned to an employee, as returned
// by the single-employee lookup endpoint.
type WorkplaceRef struct {
	// ID is the unique identifier of the workplace.
	ID int `json:"id"`
	// Name is the workplace's display name.
	Name string `json:"name"`
	// Active indicates whether the workplace is currently active.
	Active bool `json:"active"`
}

// ManagerEmployeeRef holds the subset of employee fields the API embeds for a manager.
type ManagerEmployeeRef struct {
	// ID is the unique identifier of the manager.
	ID int `json:"id"`
	// Name is the manager's full legal name.
	Name string `json:"name"`
	// AdmissionDate is the manager's hiring date in "YYYY-MM-DD" format.
	AdmissionDate string `json:"admissionDate"`
	// Fired indicates whether the manager has been terminated.
	Fired bool `json:"fired"`
	// CanViewWorkgroup indicates whether the manager has workgroup visibility permissions.
	CanViewWorkgroup bool `json:"canViewWorkgroup"`
	// Status is the numeric status code for the manager's employee record.
	Status int `json:"status"`
	// DoubleBindEmployee indicates whether the manager is shared across multiple companies.
	DoubleBindEmployee bool `json:"doubleBindEmployee"`
	// RecordsPunch indicates whether the manager uses the punch clock system.
	RecordsPunch bool `json:"recordsPunch"`
}

// ManagerRef is a reference to one of an employee's managers, as returned by the
// single-employee lookup endpoint.
type ManagerRef struct {
	// ID is the unique identifier of the manager assignment.
	ID int `json:"id"`
	// Employee holds the manager's own employee data.
	Employee ManagerEmployeeRef `json:"employee"`
}

// EmployeeDetail represents a single employee record as returned by the
// single-employee lookup endpoint. Its shape differs from Employee (returned by
// List): dates are plain "YYYY-MM-DD" strings rather than Unix timestamps, and
// several nested references carry additional fields.
type EmployeeDetail struct {
	// ID is the unique identifier of the employee.
	ID int `json:"id"`
	// Name is the employee's full legal name.
	Name string `json:"name"`
	// SocialName is the employee's preferred or social name, if provided.
	SocialName string `json:"socialName"`
	// Email is the employee's contact email address.
	Email string `json:"email"`
	// Phone is the employee's contact phone number.
	Phone string `json:"phone"`
	// CPF is the employee's Brazilian tax identification number.
	CPF string `json:"cpf"`
	// CTPS is the employee's Brazilian labor card (Carteira de Trabalho) number.
	CTPS string `json:"ctps"`
	// Series is the series number associated with the employee's CTPS.
	Series string `json:"series"`
	// PIS is the employee's Social Integration Program number, if provided.
	PIS string `json:"pis"`
	// Gender is the employee's gender as reported by the API (e.g. "MASCULINO", "FEMININO").
	Gender string `json:"gender"`
	// BirthDate is the employee's date of birth in "YYYY-MM-DD" format.
	// It is nil when the value is not present in the API response.
	BirthDate *string `json:"birthDate"`
	// AdmissionDate is the employee's hiring date in "YYYY-MM-DD" format.
	AdmissionDate string `json:"admissionDate"`
	// EffectiveDate is the date the current record became effective, as a Unix millisecond timestamp.
	EffectiveDate UnixMilliTime `json:"effectiveDate"`
	// ExternalID is an optional identifier assigned by an external system.
	ExternalID string `json:"externalId"`
	// State is the Brazilian state associated with the employee (e.g. "Pará").
	State string `json:"state"`
	// CurrentWorkSchedule is a reference to the work schedule currently active for the employee.
	CurrentWorkSchedule WorkScheduleRef `json:"currentWorkSchedule"`
	// Company is a reference to the company the employee belongs to.
	Company CompanyRef `json:"company"`
	// JobRole is a reference to the employee's current job role.
	JobRole JobRoleRef `json:"jobRoleDTO"`
	// LastManager is a reference to the employee's most recent manager.
	LastManager EntityRef `json:"lastManager"`
	// Managers holds references to all current managers for the employee.
	Managers []ManagerRef `json:"managers"`
	// WorkplaceList holds references to all workplaces assigned to the employee.
	WorkplaceList []WorkplaceRef `json:"workplaceList"`
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

// FindEmployeeParams holds optional filter parameters for the single-employee lookup endpoint.
// All fields are optional; zero values are omitted from the request.
type FindEmployeeParams struct {
	// ExternalID filters by the employee's external identifier.
	ExternalID string
	// IgnoreFired, when true, excludes terminated employees from the lookup.
	IgnoreFired bool
	// TangerinoID filters by the employee's Tangerino identifier.
	TangerinoID int
}

// Find retrieves a single employee matching the given parameters.
//
// GET /employee/find
func (s *EmployeesService) Find(ctx context.Context, params FindEmployeeParams) (*EmployeeDetail, error) {
	q := url.Values{}

	if params.ExternalID != "" {
		q.Set("externalId", params.ExternalID)
	}
	if params.IgnoreFired {
		q.Set("ignoreFired", strconv.FormatBool(params.IgnoreFired))
	}
	if params.TangerinoID != 0 {
		q.Set("tangerinoId", strconv.Itoa(params.TangerinoID))
	}

	rawURL := s.client.resolveEmployerURL("/employee/find")
	if len(q) > 0 {
		rawURL += "?" + q.Encode()
	}

	var employee EmployeeDetail
	if err := s.client.get(ctx, rawURL, &employee); err != nil {
		return nil, err
	}

	return &employee, nil
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
	if params.Size != 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}
	if params.Offset != 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.ShowFired != 0 {
		q.Set("showFired", strconv.Itoa(params.ShowFired))
	}

	rawURL := s.client.resolveEmployerURL("/employee/find-all")
	if len(q) > 0 {
		rawURL += "?" + q.Encode()
	}

	var page Page[Employee]
	if err := s.client.get(ctx, rawURL, &page); err != nil {
		return nil, err
	}

	return &page, nil
}
