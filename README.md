# Unofficial Tangerino Employer API

[![Go Reference](https://pkg.go.dev/badge/github.com/raykavin/tangerino-go.svg)](https://pkg.go.dev/github.com/raykavin/tangerino-go)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/raykavin/tangerino-go)](https://goreportcard.com/report/github.com/raykavin/tangerino-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Go client library for the **Tangerino Employer API**.

Covers the core employer modules: Employees, Companies, Holiday Calendars, Work Schedules, Workplaces, and Punches.

---

## Installation

```bash
go get github.com/raykavin/tangerino-go
```

Requires **Go 1.22+**. No external dependencies — uses only the Go standard library.

---

## Quick Start

```go
import tangerino "github.com/raykavin/tangerino-go"

client, err := tangerino.NewClient("your-username", "your-password")
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

// List employees (first page)
page, err := client.Employees.List(ctx, tangerino.ListEmployeesParams{Size: 20})
if err != nil {
    log.Fatal(err)
}

for _, e := range page.Content {
    fmt.Println(e.ID, e.Name, e.AdmissionDate.Format("02/01/2006"))
}
```

---

## Configuration

```go
// Production environment (default)
client, err := tangerino.NewClient("username", "password")

// Custom base URL (e.g. for testing)
opt, err := tangerino.WithBaseURL("https://custom.tangerino.example.com")
if err != nil {
    log.Fatal(err)
}
client, err = tangerino.NewClient("username", "password", opt)

// Staging environment
client, err = tangerino.NewClient("username", "password", tangerino.WithStagingEnv())

// Custom HTTP client (for TLS, proxies, custom timeouts, etc.)
httpClient := &http.Client{Timeout: 60 * time.Second}
client, err = tangerino.NewClient("username", "password", tangerino.WithHTTPClient(httpClient))
```

---

## Authentication

The Tangerino API uses **HTTP Basic Authentication**. Credentials are encoded and
sent automatically on every request — no additional setup is required after creating
the client.

```go
client, err := tangerino.NewClient("your-username", "your-password")
```

---

## Modules

### Employees

```go
// List employees (paginated)
page, err := client.Employees.List(ctx, tangerino.ListEmployeesParams{
    Size: 20,
})
fmt.Printf("Page 1 of %d (%d total)\n", page.TotalPages, page.TotalElements)

// Apply filters
page, err = client.Employees.List(ctx, tangerino.ListEmployeesParams{
    Size:              20,
    BranchExternalID:  "branch-001",
    ManagerExternalID: "manager-042",
    ShowFired:         1, // include terminated employees
})

// Filter by last update (Unix timestamp in milliseconds)
page, err = client.Employees.List(ctx, tangerino.ListEmployeesParams{
    LastUpdate: time.Now().Add(-24 * time.Hour).UnixMilli(),
    Size:       50,
})

// Iterate all pages
params := tangerino.ListEmployeesParams{Size: 50}
for {
    page, err := client.Employees.List(ctx, params)
    if err != nil {
        log.Fatal(err)
    }
    for _, e := range page.Content {
        fmt.Println(e.ID, e.Name)
    }
    if !page.HasNext() {
        break
    }
    params.Page = page.NextPageNumber()
}

// Access typed date fields
for _, e := range page.Content {
    fmt.Println(e.AdmissionDate.Format("02/01/2006")) // "01/04/2025"
    fmt.Println(e.AdmissionDate.Raw())                // 1743462000000
    fmt.Println(e.AdmissionDate.Time())               // time.Time

    if e.BirthDate != nil {
        fmt.Println(e.BirthDate.Format("02/01/2006"))
    }
}
```

---

### Companies

```go
// List companies (paginated)
page, err := client.Companies.List(ctx, tangerino.ListCompaniesParams{
    Size: 20,
})

// Explicit page navigation
page, err = client.Companies.List(ctx, tangerino.ListCompaniesParams{
    Page: 1,
    Size: 10,
})

for _, c := range page.Content {
    fmt.Println(c.ID, c.CNPJ, c.SocialReason, c.FantasyName)
}
```

---

### Holiday Calendars

```go
// List all holiday calendars for the employer
calendars, err := client.HolidayCalendars.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, cal := range calendars {
    fmt.Printf("[%d] %s (%d)\n", cal.ID, cal.Name, cal.Year)
    for _, h := range cal.Holidays {
        fmt.Printf("  %s - %s\n", h.Date, h.Description)
    }
}
```

---

### Work Schedules

```go
// List all work schedules
page, err := client.WorkSchedules.List(ctx, tangerino.ListWorkSchedulesParams{
    Size: 50,
})
if err != nil {
    log.Fatal(err)
}

for _, ws := range page.Content {
    fmt.Printf("[%d] %s (standard: %v, inactive: %v)\n",
        ws.ID, ws.Name, ws.Standard, ws.Inactive)

    for _, tt := range ws.Timetable {
        fmt.Printf("  Day %d: %s - %s",
            tt.Day,
            tt.StartShift1.String(), // "08:00"
            tt.EndShift1,            // *DayOffset, check nil before use
        )
        if tt.StartShift2 != nil {
            fmt.Printf(" | %s - %s",
                tt.StartShift2.String(),
                tt.EndShift2,
            )
        }
        fmt.Println()
    }

    fmt.Println("Last modified:", ws.AlterationDate.Format("02/01/2006"))
}
```

---

### Workplaces

```go
// List workplaces (paginated)
page, err := client.Workplaces.List(ctx, tangerino.ListWorkplacesParams{
    Size: 20,
})
if err != nil {
    log.Fatal(err)
}

for _, w := range page.Content {
    fmt.Printf("[%d] %s — %s, %s\n", w.ID, w.Name, w.City, w.State)
}

// Iterate all pages
params := tangerino.ListWorkplacesParams{Size: 50}
for {
    page, err := client.Workplaces.List(ctx, params)
    if err != nil {
        log.Fatal(err)
    }
    for _, w := range page.Content {
        fmt.Println(w.ID, w.Name)
    }
    if !page.HasNext() {
        break
    }
    params.Page = page.NextPageNumber()
}
```

---

### Punches

The punch endpoints live on a separate host (`apis.tangerino.com.br`) — the client handles routing automatically.

```go
adj := true
pending := false

// List punch records for an employee
punches, err := client.Punches.List(ctx, employeeID, tangerino.PunchesParams{
    Status:     3,
    Adjustment: &adj,
    StartDate:  time.Now().AddDate(0, -1, 0), // converted to Unix seconds
    EndDate:    time.Now(),
    Pending:    &pending,
})
if err != nil {
    log.Fatal(err)
}

for _, p := range punches {
    end := "open"
    if p.EndDate != nil {
        end = p.EndDate.String() // "2026-06-01T18:02:00"
    }
    fmt.Printf("[%d] %s  %s → %s  (%s manual: start=%v end=%v)\n",
        p.ID, p.Date,
        p.StartDate.String(), end,
        formatMillis(p.TotalHours),
        p.StartManual, p.EndManual,
    )
}

// Get aggregate summary for an employee
summary, err := client.Punches.Summary(ctx, employeeID, tangerino.PunchesParams{
    StartDate: time.Now().AddDate(0, -1, 0),
    EndDate:   time.Now(),
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Worked: %dms  Expected: %dms  Balance: %dms\n",
    summary.TotalWorked, summary.TotalExpected, summary.Balance)
```

**`StartDate` and `EndDate`** are `time.Time` values and are automatically converted to Unix second timestamps in the request query string.

**`Adjustment` and `Pending`** are `*bool` — use a pointer so that `false` is distinguishable from "not set":

```go
t, f := true, false
tangerino.PunchesParams{Adjustment: &t, Pending: &f}
```

---

## Custom Types

### `UnixMilliTime`

Used for absolute timestamps (`AdmissionDate`, `BirthDate`, `EffectiveDate`,
`AlterationDate`, `StartDate`). Preserves the original Unix millisecond value and
provides conversion helpers.

```go
ts := employee.AdmissionDate

ts.Raw()                     // int64   → 1776049200000  (original API value)
ts.Time()                    // time.Time (UTC)
ts.Format("02/01/2006")      // string  → "01/04/2026"
ts.Format("15:04")           // string  → "03:00" (UTC hour)
ts.Format(time.RFC3339)      // string  → "2026-04-01T03:00:00Z"
ts.String()                  // string  → "2026-04-01 03:00:00 UTC"
```

### `DayOffset`

Used for time-of-day offsets in work schedule timetables (`StartShift1`,
`EndShift1`, `StartShift2`, `EndShift2`, `StartMainInterval`, `EndMainInterval`).
Values are milliseconds from midnight and may exceed 86400000 (24 h) for shifts
that extend into the next day.

```go
d := timetable.StartShift1

d.Raw()       // int64         → 39600000    (original API value)
d.Duration()  // time.Duration → 11h0m0s
d.String()    // string        → "11:00"

// Shifts past midnight (e.g. 12x36 schedule):
// EndShift2 = 100800000 → "28:00"
```

### `LocalDateTime`

Used for punch timestamps (`StartDate`, `EndDate` in `Punch`). Represents a
wall-clock datetime without timezone in the format `"2006-01-02T15:04:05"`.

```go
p := punches[0]

p.StartDate.String()  // "2026-06-01T14:06:00"
p.StartDate.Time()    // time.Time (no timezone — server local time)

// EndDate is *LocalDateTime, nil when the employee hasn't clocked out yet
if p.EndDate != nil {
    fmt.Println(p.EndDate.String())
}
```

---

## Pagination

Paginated endpoints return `*Page[T]`, a generic type wrapping the Spring-style
`content` envelope.

```go
type Page[T any] struct {
    Content          []T
    First            bool
    Last             bool
    TotalElements    int
    TotalPages       int
    NumberOfElements int
    Size             int
    Number           int
}
```

Helper methods:

```go
page.HasNext()         // bool: whether a next page exists
page.NextPageNumber()  // int: page number to use in the next request (-1 on last page)
```

Pagination parameters follow a consistent naming convention across all services:

| Field  | Query param | Description                    |
|--------|-------------|--------------------------------|
| `Page` | `page`      | Zero-based page index          |
| `Size` | `size`      | Number of items per page       |

Full iteration pattern:

```go
params := tangerino.ListEmployeesParams{Size: 50}
for {
    page, err := client.Employees.List(ctx, params)
    if err != nil {
        log.Fatal(err)
    }
    // process page.Content ...
    if !page.HasNext() {
        break
    }
    params.Page = page.NextPageNumber()
}
```

---

## Error Handling

All methods return `*APIError` on HTTP-level failures.

```go
page, err := client.Employees.List(ctx, tangerino.ListEmployeesParams{})
if err != nil {
    switch {
    case tangerino.IsUnauthorized(err):
        // HTTP 401 — invalid credentials
    case tangerino.IsForbidden(err):
        // HTTP 403 — insufficient permissions
    case tangerino.IsNotFound(err):
        // HTTP 404 — resource not found
    case tangerino.IsRateLimited(err):
        // HTTP 429 — too many requests, back off and retry
    case tangerino.IsServerError(err):
        // HTTP 5xx — transient server error
    default:
        if apiErr, ok := err.(*tangerino.APIError); ok {
            fmt.Printf("Status: %d\n", apiErr.StatusCode)
            fmt.Printf("Body:   %s\n", apiErr.Body)
        }
    }
}
```

---

## Running Tests

```bash
go test ./... -v
```

All tests use `net/http/httptest` — no external services or environment variables
are required.

---

## Endpoints Coverage

### Employees
- [x] `GET /employee/find-all` — List employees (paginated, with filters)

### Companies
- [x] `GET /companies` — List companies (paginated)

### Holiday Calendars
- [x] `GET /holiday-calendar/` — List holiday calendars

### Work Schedules
- [x] `GET /work-schedule` — List work schedules (paginated)

### Workplaces
- [x] `GET /workplace/find-all` — List workplaces (paginated)

### Punches
- [x] `GET /punch/v2/punches/employees/{id}` — List punch records for an employee
- [x] `GET /punch/v2/punches/employees/{id}/summary` — Get aggregate time summary for an employee

**Total: 7 endpoints covered**

---

## Contributing

Contributions to tangerino-go are welcome! Here are some ways you can help:

- **Report bugs and suggest features** by opening issues on GitHub
- **Submit pull requests** with bug fixes or new features
- **Improve documentation** to help other users and developers

---

## License

tangerino-go is distributed under the **MIT License**.  
For complete license terms and conditions, see the [LICENSE](LICENSE) file in the repository.

---

## Contact

For support, collaboration, or questions about tangerino-go:

**Email**: [raykavin.meireles@gmail.com](mailto:raykavin.meireles@gmail.com)  
**GitHub**: [@raykavin](https://github.com/raykavin)
