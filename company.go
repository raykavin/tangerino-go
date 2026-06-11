package tangerino

import (
	"context"
	"net/url"
	"strconv"
)

// CompaniesService handles communication with the company endpoints.
type CompaniesService struct {
	client *Client
}

// Company represents a single company record returned by the API.
type Company struct {
	// ID is the unique identifier of the company.
	ID int `json:"id"`
	// CNPJ is the company's Brazilian federal tax registration number, formatted with punctuation.
	CNPJ string `json:"cnpj"`
	// ExternalID is an optional identifier assigned by an external system.
	ExternalID string `json:"externalId"`
	// SocialReason is the company's registered legal name.
	SocialReason string `json:"socialReason"`
	// FantasyName is the company's trade name used in day-to-day operations.
	FantasyName string `json:"fantasyName"`
	// DescriptionName is the display name used in the Tangerino interface.
	DescriptionName string `json:"descriptionName"`
}

// ListCompaniesParams holds optional filter and pagination parameters for the companies endpoint.
// All fields are optional; zero values are omitted from the request.
type ListCompaniesParams struct {
	// Page is the zero-based page index to retrieve.
	Page int
	// Size is the number of items per page.
	Size int
	// Offset is the item offset within the result set.
	Offset int
}

// List retrieves a single page of companies matching the given parameters.
// All parameters are optional; omit them by using a zero-value ListCompaniesParams.
//
// GET /companies
func (s *CompaniesService) List(ctx context.Context, params ListCompaniesParams) (*Page[Company], error) {
	q := url.Values{}

	if params.Page != 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size != 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}
	if params.Offset != 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	rawURL := s.client.resolveEmployerURL("/companies")
	if len(q) > 0 {
		rawURL += "?" + q.Encode()
	}

	var page Page[Company]
	if err := s.client.get(ctx, rawURL, &page); err != nil {
		return nil, err
	}

	return &page, nil
}
