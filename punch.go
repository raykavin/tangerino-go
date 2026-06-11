package tangerino

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// PunchesService handles communication with the punch clock endpoints.
type PunchesService struct {
	client *Client
}

const localDateTimeLayout = "2006-01-02T15:04:05"

// LocalDateTime is a wall-clock datetime without timezone, stored in the format
// "2006-01-02T15:04:05" as returned by the punch API.
// Use Time() to obtain a time.Time (interpreted in the local timezone of the server).
type LocalDateTime struct {
	t time.Time
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *LocalDateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse(localDateTimeLayout, s)
	if err != nil {
		return fmt.Errorf("tangerino: parsing LocalDateTime %q: %w", s, err)
	}
	d.t = t
	return nil
}

// MarshalJSON implements json.Marshaler.
func (d LocalDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.t.Format(localDateTimeLayout))
}

// Time returns the underlying time.Time value.
func (d LocalDateTime) Time() time.Time { return d.t }

// String returns the value formatted as "2006-01-02T15:04:05".
func (d LocalDateTime) String() string { return d.t.Format(localDateTimeLayout) }

// Punch represents a single work interval record (a pair of clock-in / clock-out events).
// Each punch covers one continuous work block; a day with a lunch break typically has two.
type Punch struct {
	// ID is the unique identifier of the punch record.
	ID int `json:"id"`
	// Date is the calendar date of the punch in "YYYY-MM-DD" format.
	Date string `json:"date"`
	// StartDate is the clock-in time for this interval.
	StartDate LocalDateTime `json:"startDate"`
	// EndDate is the clock-out time for this interval.
	// It is nil when the employee has not yet clocked out.
	EndDate *LocalDateTime `json:"endDate"`
	// Status is the numeric status code of the punch record (e.g. 2 = approved).
	Status int `json:"status"`
	// StartManual indicates whether the clock-in was entered manually.
	StartManual bool `json:"startManual"`
	// EndManual indicates whether the clock-out was entered manually.
	EndManual bool `json:"endManual"`
	// Pending indicates whether the punch is awaiting approval.
	Pending bool `json:"pending"`
	// Accredited indicates whether the punch has been credited to the hours bank.
	Accredited bool `json:"accredited"`
	// Adjustment indicates whether this punch record was adjusted.
	Adjustment bool `json:"adjustment"`
	// Canceled indicates whether the punch has been voided.
	Canceled bool `json:"canceled"`
	// TotalHours is the duration of the interval in milliseconds.
	// Despite the name, the API returns a millisecond value (e.g. 14400000 = 4 h).
	TotalHours int64 `json:"totalHours"`
}

// PunchSummary holds aggregate time data for an employee over a queried period.
type PunchSummary struct {
	// EmployeeID is the identifier of the employee.
	EmployeeID int `json:"employeeId"`
	// TotalWorked is the total time worked in milliseconds.
	TotalWorked int64 `json:"totalWorked"`
	// TotalExpected is the expected work time in milliseconds.
	TotalExpected int64 `json:"totalExpected"`
	// Balance is the difference between worked and expected time in milliseconds.
	// Positive values indicate overtime; negative values indicate a deficit.
	Balance int64 `json:"balance"`
	// Absences is the number of absence days in the period.
	Absences int `json:"absences"`
	// Delays is the total delay time in milliseconds.
	Delays int64 `json:"delays"`
	// Overtime is the total overtime in milliseconds.
	Overtime int64 `json:"overtime"`
	// HoursBank is the accumulated hours-bank balance in milliseconds.
	HoursBank int64 `json:"hoursBank"`
}

// PunchesParams holds filter parameters for punch clock endpoints.
// StartDate and EndDate accept time.Time values and are sent to the API as
// Unix second timestamps. Zero-value time.Time fields are omitted.
// Adjustment and Pending are pointer booleans so that false can be distinguished
// from "not set".
type PunchesParams struct {
	// Status filters punches by status code (e.g. 3 for pending approval).
	Status int
	// Adjustment filters by whether the punch was manually adjusted.
	// Pass a pointer to true or false; nil omits the parameter.
	Adjustment *bool
	// StartDate is the inclusive start of the date range.
	// Converted to a Unix second timestamp when building the request.
	StartDate time.Time
	// EndDate is the inclusive end of the date range.
	// Converted to a Unix second timestamp when building the request.
	EndDate time.Time
	// Pending filters by whether the punch is awaiting approval.
	// Pass a pointer to true or false; nil omits the parameter.
	Pending *bool
}

func buildPunchQuery(params PunchesParams) url.Values {
	q := url.Values{}

	if params.Status != 0 {
		q.Set("status", strconv.Itoa(params.Status))
	}
	if params.Adjustment != nil {
		q.Set("adjustment", strconv.FormatBool(*params.Adjustment))
	}
	if !params.StartDate.IsZero() {
		q.Set("startDate", strconv.FormatInt(params.StartDate.Unix(), 10))
	}
	if !params.EndDate.IsZero() {
		q.Set("endDate", strconv.FormatInt(params.EndDate.Unix(), 10))
	}
	if params.Pending != nil {
		q.Set("pending", strconv.FormatBool(*params.Pending))
	}

	return q
}

// List retrieves all punch records for the given employee matching the filter params.
//
// GET /punch/v2/punches/employees/{employeeID}
func (s *PunchesService) List(ctx context.Context, employeeID int, params PunchesParams) ([]Punch, error) {
	q := buildPunchQuery(params)

	rawURL := s.client.resolvePunchURL(fmt.Sprintf("/punch/v2/punches/employees/%d", employeeID))
	if len(q) > 0 {
		rawURL += "?" + q.Encode()
	}

	var punches []Punch
	if err := s.client.get(ctx, rawURL, &punches); err != nil {
		return nil, err
	}

	return punches, nil
}
