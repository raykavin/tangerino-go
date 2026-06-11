package tangerino

import "context"

// HolidayCalendarsService handles communication with the holiday calendar endpoints.
type HolidayCalendarsService struct {
	client *Client
}

// Holiday represents a single public or regional holiday entry within a calendar.
type Holiday struct {
	// ID is the unique identifier of the holiday.
	ID int `json:"id"`
	// Description is the name or label of the holiday.
	Description string `json:"description"`
	// Date is the calendar date of the holiday in YYYY-MM-DD format.
	Date string `json:"date"`
}

// HolidayCalendar represents a named collection of holidays for a specific year.
// Each calendar may cover a national, regional, or custom set of holidays.
type HolidayCalendar struct {
	// ID is the unique identifier of the calendar.
	ID int `json:"id"`
	// Name is the human-readable label assigned to the calendar.
	Name string `json:"name"`
	// Description is an optional extended description of the calendar.
	Description string `json:"description"`
	// Year is the year this calendar applies to.
	Year int `json:"year"`
	// Holidays lists the individual holiday entries contained in this calendar.
	Holidays []Holiday `json:"holidays"`
}

// holidayCalendarListResponse is the response envelope returned by the holiday-calendar endpoint.
type holidayCalendarListResponse struct {
	Code     int               `json:"code"`
	Status   string            `json:"status"`
	Location string            `json:"location"`
	Messages []string          `json:"messages"`
	Item     []HolidayCalendar `json:"item"`
}

// List retrieves all holiday calendars available for the authenticated employer.
//
// GET /holiday-calendar/
func (s *HolidayCalendarsService) List(ctx context.Context) ([]HolidayCalendar, error) {
	rawURL := s.client.resolveURL("/holiday-calendar/")

	var resp holidayCalendarListResponse
	if err := s.client.get(ctx, rawURL, &resp); err != nil {
		return nil, err
	}

	return resp.Item, nil
}
