package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
		PageSize: 20,
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

var _ = allEmployees
var _ = os.Args
