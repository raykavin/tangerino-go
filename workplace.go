package tangerino

import (
	"context"
	"net/url"
	"strconv"
)

// WorkplacesService handles communication with the workplace endpoints.
type WorkplacesService struct {
	client *Client
}

// Workplace represents a single workplace record returned by the API.
type Workplace struct {
	// ID is the unique identifier of the workplace.
	ID int `json:"id"`
	// Name is the human-readable label for the workplace.
	Name string `json:"name"`
	// ExternalID is an optional identifier assigned by an external system.
	ExternalID string `json:"externalId"`
	// Company is a reference to the company this workplace belongs to.
	Company EntityRef `json:"company"`
	// Address is the street address of the workplace.
	Address string `json:"address"`
	// City is the city where the workplace is located.
	City string `json:"city"`
	// State is the state or province where the workplace is located.
	State string `json:"state"`
	// ZipCode is the postal code of the workplace.
	ZipCode string `json:"zipCode"`
	// Inactive indicates whether the workplace has been deactivated.
	Inactive bool `json:"inactive"`
}

// ListWorkplacesParams holds optional pagination parameters for the workplace list endpoint.
// All fields are optional; zero values are omitted from the request.
type ListWorkplacesParams struct {
	// Page is the zero-based page index to retrieve.
	Page int
	// Size is the number of items per page.
	Size int
}

// List retrieves a single page of workplaces for the authenticated employer.
// All parameters are optional; omit them by using a zero-value ListWorkplacesParams.
//
// GET /workplace/find-all
func (s *WorkplacesService) List(ctx context.Context, params ListWorkplacesParams) (*Page[Workplace], error) {
	q := url.Values{}

	if params.Page != 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size != 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}

	rawURL := s.client.resolveEmployerURL("/workplace/find-all")
	if len(q) > 0 {
		rawURL += "?" + q.Encode()
	}

	var page Page[Workplace]
	if err := s.client.get(ctx, rawURL, &page); err != nil {
		return nil, err
	}

	return &page, nil
}
