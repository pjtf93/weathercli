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

func TestCurrentWeatherE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	weather, err := client.Current(ctx, "Paris")
	if err != nil {
		t.Fatalf("Current failed: %v", err)
	}

	if weather.Location.Name == "" {
		t.Error("Location name is empty")
	}
	if weather.Time.IsZero() {
		t.Error("Weather time is zero")
	}
	if weather.Condition == "" {
		t.Error("Weather condition is empty")
	}
	if weather.Temperature < -100 || weather.Temperature > 100 {
		t.Errorf("Temperature out of reasonable range: %.1f", weather.Temperature)
	}
	if weather.Humidity < 0 || weather.Humidity > 100 {
		t.Errorf("Humidity out of range: %d", weather.Humidity)
	}
}

func TestCurrentWeatherByCoordsE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	// London coordinates
	lat := 51.5074
	lon := -0.1278

	weather, err := client.CurrentByCoords(ctx, lat, lon, nil)
	if err != nil {
		t.Fatalf("CurrentByCoords failed: %v", err)
	}

	if weather.Time.IsZero() {
		t.Error("Weather time is zero")
	}
	
	// API may round coordinates slightly, allow small difference
	latDiff := weather.Location.Latitude - lat
	if latDiff < -0.01 || latDiff > 0.01 {
		t.Errorf("Latitude too different: got %.4f, want %.4f (diff %.4f)", weather.Location.Latitude, lat, latDiff)
	}
	lonDiff := weather.Location.Longitude - lon
	if lonDiff < -0.01 || lonDiff > 0.01 {
		t.Errorf("Longitude too different: got %.4f, want %.4f (diff %.4f)", weather.Location.Longitude, lon, lonDiff)
	}
}

func TestForecastE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	// Test daily forecast
	forecast, err := client.Forecast(ctx, "Berlin", 3, false)
	if err != nil {
		t.Fatalf("Forecast failed: %v", err)
	}

	if len(forecast.Daily) != 3 {
		t.Errorf("Expected 3 days, got %d", len(forecast.Daily))
	}

	for i, day := range forecast.Daily {
		if day.Date.IsZero() {
			t.Errorf("Day %d: date is zero", i)
		}
		if day.Sunrise.IsZero() {
			t.Errorf("Day %d: sunrise is zero", i)
		}
		if day.Sunset.IsZero() {
			t.Errorf("Day %d: sunset is zero", i)
		}
		if day.Condition == "" {
			t.Errorf("Day %d: condition is empty", i)
		}
		if day.TempMax < day.TempMin {
			t.Errorf("Day %d: max temp (%.1f) < min temp (%.1f)", i, day.TempMax, day.TempMin)
		}
	}
}

func TestForecastHourlyE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	// Test hourly forecast
	forecast, err := client.Forecast(ctx, "Tokyo", 1, true)
	if err != nil {
		t.Fatalf("Forecast (hourly) failed: %v", err)
	}

	if len(forecast.Hourly) == 0 {
		t.Error("Expected hourly data")
	}

	for i, hour := range forecast.Hourly {
		if hour.Time.IsZero() {
			t.Errorf("Hour %d: time is zero", i)
		}
		if hour.Condition == "" {
			t.Errorf("Hour %d: condition is empty", i)
		}
		if hour.Temperature < -100 || hour.Temperature > 100 {
			t.Errorf("Hour %d: temperature out of range: %.1f", i, hour.Temperature)
		}
	}
}

func TestForecastByCoordsE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	ctx := context.Background()

	// New York coordinates
	lat := 40.7128
	lon := -74.0060

	forecast, err := client.ForecastByCoords(ctx, lat, lon, 2, false, nil)
	if err != nil {
		t.Fatalf("ForecastByCoords failed: %v", err)
	}

	if len(forecast.Daily) != 2 {
		t.Errorf("Expected 2 days, got %d", len(forecast.Daily))
	}
}

func TestLocationFields(t *testing.T) {
	loc := Location{
		Name:      "London",
		Latitude:  51.5074,
		Longitude: -0.1278,
		Country:   "United Kingdom",
		Admin1:    "England",
		Timezone:  "Europe/London",
	}

	if loc.Name != "London" {
		t.Errorf("Name = %q, want %q", loc.Name, "London")
	}
	if loc.Country != "United Kingdom" {
		t.Errorf("Country = %q, want %q", loc.Country, "United Kingdom")
	}
}

func TestWeatherCodeEdgeCases(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{-1, "Unknown"},
		{0, "Clear sky"},
		{3, "Overcast"},
		{45, "Foggy"},
		{51, "Light drizzle"},
		{61, "Slight rain"},
		{71, "Slight snow"},
		{80, "Slight rain showers"},
		{95, "Thunderstorm"},
		{99, "Thunderstorm with heavy hail"},
		{100, "Unknown"},
		{1000, "Unknown"},
	}

	for _, tt := range tests {
		got := GetCondition(tt.code)
		if got != tt.want {
			t.Errorf("GetCondition(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestWindDirectionEdgeCases(t *testing.T) {
	tests := []struct {
		degrees int
		want    string
	}{
		{-45, "NW"},   // -45 + 360 = 315 = NW
		{-1, "N"},     // -1 + 360 = 359 = N
		{22, "NNE"},   // 22 = NNE
		{23, "NNE"},   // 23 = NNE
		{34, "NE"},    // 34 = NE
		{67, "ENE"},   // 67 = ENE
		{90, "E"},     // 90 = E
		{361, "N"},    // 361 % 360 = 1 = N
		{405, "NE"},   // 405 % 360 = 45 = NE
		{-90, "W"},    // -90 + 360 = 270 = W
		{-180, "S"},   // -180 + 360 = 180 = S
	}

	for _, tt := range tests {
		got := WindDirection(tt.degrees)
		if got != tt.want {
			t.Errorf("WindDirection(%d) = %q, want %q", tt.degrees, got, tt.want)
		}
	}
}

func TestClientOptions(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "empty options",
			opts: Options{},
		},
		{
			name: "custom baseURL only",
			opts: Options{BaseURL: "https://custom.api"},
		},
		{
			name: "custom geoBaseURL only",
			opts: Options{GeoBaseURL: "https://geo.api"},
		},
		{
			name: "custom timeout only",
			opts: Options{Timeout: 5},
		},
		{
			name: "all custom",
			opts: Options{
				BaseURL:    "https://custom.api",
				GeoBaseURL: "https://geo.api",
				Timeout:    5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.opts)
			if client == nil {
				t.Fatal("NewClient returned nil")
			}
			if client.httpClient == nil {
				t.Error("httpClient is nil")
			}
		})
	}
}
