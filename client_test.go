package tangerino_test

import (
	"testing"

	tangerino "github.com/raykavin/tangerino-go"
)

func TestNewClient_RequiresCredentials(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"empty username", "", "pass", true},
		{"empty password", "user", "", true},
		{"both empty", "", "", true},
		{"valid credentials", "user", "pass", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tangerino.NewClient(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithBaseURL_RejectsInvalidURL(t *testing.T) {
	_, err := tangerino.WithBaseURL("://bad-url")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestWithBaseURL_AcceptsValidURL(t *testing.T) {
	opt, err := tangerino.WithBaseURL("https://custom.example.com")
	if err != nil {
		t.Fatalf("WithBaseURL: unexpected error: %v", err)
	}
	if opt == nil {
		t.Error("expected non-nil Option")
	}
}
