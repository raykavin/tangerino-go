// Package tangerino provides a client library for the Tangerino employer API.
//
// Create a client with NewClient, then use its service fields to access API resources:
//
//	client, err := tangerino.NewClient("username", "password")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	page, err := client.Employees.List(ctx, tangerino.ListEmployeesParams{PageSize: 20})
//	punches, err := client.Punches.GetEmployeePunches(ctx, employeeID)
package tangerino

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "tangerino-go/1.0.0"
	acceptHeader     = "application/json;charset=UTF-8"
)

// Client is the Tangerino API client.
// Use its service fields to call specific API resources.
type Client struct {
	username   string
	password   string
	env        Environment
	httpClient *http.Client

	// Employees provides access to employee management endpoints.
	Employees *EmployeesService
	// HolidayCalendars provides access to holiday calendar endpoints.
	HolidayCalendars *HolidayCalendarsService
	// WorkSchedules provides access to work schedule endpoints.
	WorkSchedules *WorkSchedulesService
	// Companies provides access to company endpoints.
	Companies *CompaniesService
}

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithHTTPClient replaces the default HTTP client with a custom one.
// Use this to configure custom timeouts, TLS settings, or proxies.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithBaseURL overrides the API base URL.
// Returns an error if the provided value is not a valid URL.
func WithBaseURL(rawURL string) (Option, error) {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", rawURL, err)
	}
	return func(c *Client) {
		c.env.baseURL = rawURL
	}, nil
}

// WithStagingEnv points the client at the staging environment.
func WithStagingEnv() Option {
	return func(c *Client) {
		c.env = stagingEnv
	}
}

// NewClient creates an authenticated Tangerino API client.
// Both username and password are required for Basic Authentication and
// are used on every outgoing request.
func NewClient(username, password string, opts ...Option) (*Client, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	c := &Client{
		username: username,
		password: password,
		env:      prodEnv,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Employees = &EmployeesService{client: c}
	c.HolidayCalendars = &HolidayCalendarsService{client: c}
	c.WorkSchedules = &WorkSchedulesService{client: c}
	c.Companies = &CompaniesService{client: c}

	return c, nil
}

// basicAuth returns the value for the Authorization header using HTTP Basic Authentication.
func (c *Client) basicAuth() string {
	encoded := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
	return "Basic " + encoded
}

// resolveURL builds the full URL by appending path to the configured base URL.
func (c *Client) resolveURL(path string) string {
	return strings.TrimRight(c.env.baseURL, "/") + path
}

// get performs an authenticated GET request against rawURL and decodes the
// JSON response body into result. It returns an APIError for any HTTP 4xx or
// 5xx response.
func (c *Client) get(ctx context.Context, rawURL string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", acceptHeader)
	req.Header.Set("Authorization", c.basicAuth())
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode, Body: body}
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}
