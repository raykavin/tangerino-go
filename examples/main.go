package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tangerino "github.com/raykavin/tangerino-go"
)

func main() {
	username := os.Getenv("TANGERINO_USERNAME")
	password := os.Getenv("TANGERINO_PASSWORD")

	client, err := tangerino.NewClient(username, password)
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}

	ctx := context.Background()

	listEmployees(ctx, client)
}

func listEmployees(ctx context.Context, client *tangerino.Client) {
	page, err := client.Employees.List(ctx, tangerino.ListEmployeesParams{
		Size: 20,
	})
	if err != nil {
		handleError("listing employees", err)
	}

	fmt.Printf("Employees: page %d of %d (%d total)\n\n", page.Number+1, page.TotalPages, page.TotalElements)

	for _, e := range page.Content {
		fmt.Printf("  ID: %-10d  Name: %s\n", e.ID, e.Name)
	}

	if page.HasNext() {
		fmt.Printf("\n  Next page: %d\n", page.NextPageNumber())
	}
}

func listWorkplaces(ctx context.Context, client *tangerino.Client) {
	page, err := client.Workplaces.List(ctx, tangerino.ListWorkplacesParams{
		Size: 20,
	})
	if err != nil {
		handleError("listing workplaces", err)
	}

	fmt.Printf("Workplaces: page %d of %d (%d total)\n\n", page.Number+1, page.TotalPages, page.TotalElements)

	for _, w := range page.Content {
		fmt.Printf("  ID: %-10d  Name: %s\n", w.ID, w.Name)
	}

	if page.HasNext() {
		fmt.Printf("\n  Next page: %d\n", page.NextPageNumber())
	}
}

func handleError(op string, err error) {
	switch {
	case tangerino.IsUnauthorized(err):
		log.Fatal("invalid credentials")
	case tangerino.IsNotFound(err):
		log.Fatalf("%s: not found", op)
	case tangerino.IsServerError(err):
		log.Fatalf("%s: server error: %v", op, err)
	default:
		log.Fatalf("%s: %v", op, err)
	}
}

// Demonstrates iterating all pages of employees.
func allEmployees(ctx context.Context, client *tangerino.Client) ([]tangerino.Employee, error) {
	var all []tangerino.Employee
	params := tangerino.ListEmployeesParams{Size: 20}

	for {
		page, err := client.Employees.List(ctx, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page.Content...)

		if !page.HasNext() {
			break
		}

		params.Page = page.NextPageNumber()
	}

	return all, nil
}

// Demonstrates iterating all pages of workplaces.
func allWorkplaces(ctx context.Context, client *tangerino.Client) ([]tangerino.Workplace, error) {
	var all []tangerino.Workplace
	params := tangerino.ListWorkplacesParams{Size: 20}

	for {
		page, err := client.Workplaces.List(ctx, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page.Content...)

		if !page.HasNext() {
			break
		}

		params.Page = page.NextPageNumber()
	}

	return all, nil
}

func listPunches(ctx context.Context, client *tangerino.Client, employeeID int) {
	adj := true
	pending := false

	punches, err := client.Punches.List(ctx, employeeID, tangerino.PunchesParams{
		Status:     3,
		Adjustment: &adj,
		StartDate:  time.Now().AddDate(0, -1, 0),
		EndDate:    time.Now(),
		Pending:    &pending,
	})
	if err != nil {
		handleError("listing punches", err)
	}

	fmt.Printf("Punches for employee %d: %d record(s)\n\n", employeeID, len(punches))

	for _, p := range punches {
		end := "open"
		if p.EndDate != nil {
			end = p.EndDate.String()
		}
		fmt.Printf("  ID: %-12d  %s  %s → %s\n", p.ID, p.Date, p.StartDate, end)
	}
}

var _ = allEmployees
var _ = allWorkplaces
var _ = listWorkplaces
var _ = listPunches
var _ = os.Args
