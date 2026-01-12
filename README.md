# ☀️ weathercli — Simple Weather CLI for LLMs & Humans

[![CI](https://github.com/pjtf93/weathercli/workflows/CI/badge.svg)](https://github.com/pjtf93/weathercli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/pjtf93/weathercli)](https://goreportcard.com/report/github.com/pjtf93/weathercli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Fast weather CLI built for LLM integration and human readability. Get current weather or forecasts for any location worldwide.

## Highlights

- Current weather conditions with temp, humidity, wind
- Daily forecast (up to 16 days)
- Hourly forecast (up to 16 days)
- Location search by name (geocoding)
- JSON output for LLM parsing
- Color human output (respects `NO_COLOR`)
- Free API (Open-Meteo) — no key required
- Works worldwide

## Install / Run

```bash
# From source
cd weathercli
go build -o weathercli cmd/weathercli/main.go

# Or with go install
go install github.com/pjtf93/weathercli/cmd/weathercli@latest
```

## CLI

```text
weathercli [--json] [--no-color] [--verbose] <command>

Commands:
  current   Get current weather for a location
  forecast  Get weather forecast for a location
  search    Search for location coordinates
```

### Current Weather

```bash
# By location name
weathercli current "New York"
weathercli current "London, UK"
weathercli current "Tokyo, Japan"

# JSON output for LLMs
weathercli current "Paris" --json
```

### Forecast

```bash
# Daily forecast (default: 7 days)
weathercli forecast "Berlin" --days 7

# Hourly forecast
weathercli forecast "Madrid" --hourly --hours 24

# JSON output
weathercli forecast "Sydney" --days 5 --json
```

### Search Locations

```bash
# Find coordinates for a location
weathercli search "San Francisco"
weathercli search "Barcelona" --json
```

## Library Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/pjtf93/weathercli"
)

func main() {
    client := weathercli.NewClient()
    
    // Get current weather
    weather, err := client.Current(context.Background(), "London")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Temperature: %.1f°C\n", weather.Temperature)
    
    // Get forecast
    forecast, err := client.Forecast(context.Background(), "Paris", 7)
    if err != nil {
        panic(err)
    }
    for _, day := range forecast.Daily {
        fmt.Printf("%s: %.1f°C\n", day.Date, day.TempMax)
    }
}
```

## API

Uses [Open-Meteo](https://open-meteo.com/) - a free weather API:
- No API key required
- No rate limits for non-commercial use
- Global coverage
- High accuracy data from multiple sources

## Testing

```bash
go test ./...
go test -v ./...
```

## LLM Integration

Perfect for AI agents and scripts:

```bash
# Get structured JSON output
weathercli current "San Francisco" --json | jq .temperature
weathercli forecast "Berlin" --days 5 --json | jq '.daily[0].condition'

# Example LLM tool/function call
{
  "name": "get_weather",
  "description": "Get current weather for a location",
  "parameters": {
    "location": "Paris, France"
  },
  "command": "weathercli current 'Paris, France' --json"
}
```

### JSON Output Structure

#### Current Weather
```json
{
  "location": {
    "name": "London",
    "latitude": 51.5074,
    "longitude": -0.1278,
    "country": "United Kingdom",
    "timezone": "Europe/London"
  },
  "time": "2024-01-12T14:30:00Z",
  "temperature": 10.3,
  "apparent": 8.1,
  "humidity": 86,
  "wind_speed": 12.2,
  "wind_direction": 202,
  "condition": "Overcast",
  "weather_code": 3
}
```

#### Forecast
```json
{
  "location": { ... },
  "daily": [
    {
      "date": "2024-01-12",
      "temp_max": 12.1,
      "temp_min": 4.3,
      "condition": "Slight rain",
      "precip_prob": 75,
      "sunrise": "2024-01-12T08:04:00Z",
      "sunset": "2024-01-12T16:45:00Z"
    }
  ]
}
```

## Notes

- All temperatures in Celsius
- Wind speed in km/h
- Precipitation in mm / snowfall in cm
- Times in location's local timezone
- Geocoding uses Open-Meteo's built-in service
- No API key required, free for non-commercial use
- Weather codes follow WMO standard (0-99)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Development Setup

```bash
# Clone and test
git clone https://github.com/pjtf93/weathercli.git
cd weathercli
go test ./...

# Build
make build
# or
go build -o weathercli ./cmd/weathercli
```

## License

MIT License - see [LICENSE](LICENSE) for details.
