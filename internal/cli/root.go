package cli

import "time"

// Root defines the CLI command tree.
type Root struct {
	Global   GlobalOptions `embed:""`
	Current  CurrentCmd    `cmd:"" help:"Get current weather for a location."`
	Forecast ForecastCmd   `cmd:"" help:"Get weather forecast for a location."`
	Search   SearchCmd     `cmd:"" help:"Search for location coordinates."`
}

// GlobalOptions are flags shared by all commands.
type GlobalOptions struct {
	BaseURL    string        `help:"Weather API base URL." env:"WEATHER_BASE_URL" default:"https://api.open-meteo.com/v1"`
	GeoBaseURL string        `help:"Geocoding API base URL." env:"WEATHER_GEO_BASE_URL" default:"https://geocoding-api.open-meteo.com/v1"`
	Timeout    time.Duration `help:"HTTP timeout." default:"10s"`
	JSON       bool          `help:"Output JSON."`
	NoColor    bool          `help:"Disable color output."`
	Verbose    bool          `help:"Verbose logging."`
	Version    VersionFlag   `name:"version" help:"Print version and exit."`
}

// CurrentCmd gets current weather.
type CurrentCmd struct {
	Location string `arg:"" name:"location" help:"Location name (e.g. 'New York', 'London, UK')."`
}

// ForecastCmd gets weather forecast.
type ForecastCmd struct {
	Location string `arg:"" name:"location" help:"Location name (e.g. 'Paris', 'Tokyo, Japan')."`
	Days     int    `help:"Number of forecast days (1-16)." default:"7"`
	Hourly   bool   `help:"Show hourly forecast instead of daily."`
	Hours    int    `help:"Number of hours for hourly forecast (1-384)." default:"24"`
}

// SearchCmd searches for locations.
type SearchCmd struct {
	Query string `arg:"" name:"query" help:"Location search query."`
	Limit int    `help:"Max results (1-10)." default:"5"`
}
