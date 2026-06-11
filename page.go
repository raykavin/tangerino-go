package tangerino

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
