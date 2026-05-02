package tangerino

import "context"

// WorkSchedulesService handles communication with the work schedule endpoints.
type WorkSchedulesService struct {
	client *Client
}

// WorkScheduleTimetable represents the time configuration for a single day within a work schedule.
// All interval values are milliseconds elapsed since midnight.
type WorkScheduleTimetable struct {
	// ID is the unique identifier of this timetable entry.
	ID int `json:"id"`
	// Day is the day of the week this entry applies to (1=Sunday, 2=Monday, ..., 7=Saturday).
	Day int `json:"day"`
	// StartMainInterval is the start of the main break period as a day offset.
	StartMainInterval DayOffset `json:"startMainInterval"`
	// EndMainInterval is the end of the main break period as a day offset.
	// It is nil when the schedule has no defined end for the main interval.
	EndMainInterval *DayOffset `json:"endMainInterval"`
	// StartShift1 is the start of the first shift as a day offset.
	StartShift1 DayOffset `json:"startShift1"`
	// EndShift1 is the end of the first shift as a day offset.
	// It is nil when the first shift has no defined end time.
	EndShift1 *DayOffset `json:"endShift1"`
	// StartShift2 is the start of the second shift as a day offset.
	// It is nil when there is no second shift.
	StartShift2 *DayOffset `json:"startShift2"`
	// EndShift2 is the end of the second shift as a day offset.
	// It is nil when there is no second shift.
	EndShift2 *DayOffset `json:"endShift2"`
	// IntervalPreAssigned1And2 indicates whether the interval between shifts 1 and 2 is pre-assigned.
	IntervalPreAssigned1And2 bool `json:"intervalPreAssigned1And2"`
	// IntervalPreAssigned2And3 indicates whether the interval between shifts 2 and 3 is pre-assigned.
	IntervalPreAssigned2And3 bool `json:"intervalPreAssigned2And3"`
	// IntervalPreAssigned3And4 indicates whether the interval between shifts 3 and 4 is pre-assigned.
	IntervalPreAssigned3And4 bool `json:"intervalPreAssigned3And4"`
	// IntervalPreAssigned4And5 indicates whether the interval between shifts 4 and 5 is pre-assigned.
	IntervalPreAssigned4And5 bool `json:"intervalPreAssigned4And5"`
	// IntervalPreAssigned5And6 indicates whether the interval between shifts 5 and 6 is pre-assigned.
	IntervalPreAssigned5And6 bool `json:"intervalPreAssigned5And6"`
}

// WorkSchedule represents a full work schedule definition including its daily timetables.
type WorkSchedule struct {
	// ID is the unique identifier of the work schedule.
	ID int `json:"id"`
	// Name is the human-readable label for the work schedule.
	Name string `json:"name"`
	// Standard indicates whether this is the default system schedule.
	Standard bool `json:"standard"`
	// Timetable holds the per-day time configurations for this schedule.
	Timetable []WorkScheduleTimetable `json:"workScheduleTimetableList"`
	// AlterationDate is the Unix millisecond timestamp of the last modification.
	AlterationDate UnixMilliTime `json:"alterationDate"`
	// PreAssignedInterval indicates whether break intervals are pre-assigned across the schedule.
	PreAssignedInterval bool `json:"preAssignedInterval"`
	// ShowIntradayInTimeSheet indicates whether intraday entries appear in the time sheet.
	ShowIntradayInTimeSheet bool `json:"showIntradayInTimeSheet"`
	// IgnoreHoliday indicates whether this schedule applies on public holidays.
	IgnoreHoliday bool `json:"ignoreHoliday"`
	// Inactive indicates whether the schedule has been deactivated.
	Inactive bool `json:"inactive"`
}

// List retrieves all work schedules available for the authenticated employer.
//
// GET /work-schedule
func (s *WorkSchedulesService) List(ctx context.Context) (*Page[WorkSchedule], error) {
	rawURL := s.client.resolveURL("/work-schedule")

	var page Page[WorkSchedule]
	if err := s.client.get(ctx, rawURL, &page); err != nil {
		return nil, err
	}

	return &page, nil
}
