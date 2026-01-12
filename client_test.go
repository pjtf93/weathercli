package weathercli

import (
	"context"
	"testing"
)

func TestGetCondition(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{0, "Clear sky"},
		{1, "Mainly clear"},
		{61, "Slight rain"},
		{95, "Thunderstorm"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		got := GetCondition(tt.code)
		if got != tt.want {
			t.Errorf("GetCondition(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestWindDirection(t *testing.T) {
	tests := []struct {
		degrees int
		want    string
	}{
		{0, "N"},
		{45, "NE"},
		{90, "E"},
		{135, "SE"},
		{180, "S"},
		{225, "SW"},
		{270, "W"},
		{315, "NW"},
		{360, "N"},
	}

	for _, tt := range tests {
		got := WindDirection(tt.degrees)
		if got != tt.want {
			t.Errorf("WindDirection(%d) = %q, want %q", tt.degrees, got, tt.want)
		}
	}
}

func TestClientCreation(t *testing.T) {
	// Default client
	c1 := NewClient()
	if c1 == nil {
		t.Fatal("NewClient() returned nil")
	}
	if c1.baseURL != defaultBaseURL {
		t.Errorf("Default baseURL = %q, want %q", c1.baseURL, defaultBaseURL)
	}

	// Custom options
	customBase := "https://custom.api.com"
	c2 := NewClient(Options{BaseURL: customBase})
	if c2.baseURL != customBase {
		t.Errorf("Custom baseURL = %q, want %q", c2.baseURL, customBase)
	}
}

// Integration test - only runs if WEATHER_E2E is set
func TestSearchLocationE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	locations, err := client.SearchLocation(ctx, "London")
	if err != nil {
		t.Fatalf("SearchLocation failed: %v", err)
	}

	if len(locations) == 0 {
		t.Fatal("Expected at least one location")
	}

	loc := locations[0]
	if loc.Name == "" {
		t.Error("Location name is empty")
	}
	if loc.Latitude == 0 && loc.Longitude == 0 {
		t.Error("Location coordinates are zero")
	}
}

func TestSearchLocationNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	_, err := client.SearchLocation(ctx, "ThisLocationDoesNotExist12345")
	if err == nil {
		t.Error("Expected error for non-existent location")
	}
}
