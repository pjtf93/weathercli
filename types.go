package weathercli

import "time"

// Location represents a geographic location with coordinates.
type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country,omitempty"`
	Admin1    string  `json:"admin1,omitempty"` // State/Province
	Timezone  string  `json:"timezone,omitempty"`
}

// CurrentWeather represents current weather conditions.
type CurrentWeather struct {
	Location      Location  `json:"location"`
	Time          time.Time `json:"time"`
	Temperature   float64   `json:"temperature"`   // °C
	Apparent      float64   `json:"apparent"`      // Feels like °C
	Humidity      int       `json:"humidity"`      // %
	Precipitation float64   `json:"precipitation"` // mm
	Rain          float64   `json:"rain"`          // mm
	Snowfall      float64   `json:"snowfall"`      // cm
	WindSpeed     float64   `json:"wind_speed"`    // km/h
	WindDirection int       `json:"wind_direction"`
	Pressure      float64   `json:"pressure"` // hPa
	CloudCover    int       `json:"cloud_cover"`
	Visibility    float64   `json:"visibility"` // meters
	UVIndex       float64   `json:"uv_index"`
	WeatherCode   int       `json:"weather_code"`
	Condition     string    `json:"condition"` // Human-readable
}

// DailyForecast represents a single day's forecast.
type DailyForecast struct {
	Date          time.Time `json:"date"`
	TempMax       float64   `json:"temp_max"`
	TempMin       float64   `json:"temp_min"`
	ApparentMax   float64   `json:"apparent_max"`
	ApparentMin   float64   `json:"apparent_min"`
	Precipitation float64   `json:"precipitation"`
	Rain          float64   `json:"rain"`
	Snowfall      float64   `json:"snowfall"`
	WindSpeedMax  float64   `json:"wind_speed_max"`
	WindDirection int       `json:"wind_direction"`
	UVIndexMax    float64   `json:"uv_index_max"`
	PrecipProb    int       `json:"precip_prob"` // %
	Sunrise       time.Time `json:"sunrise"`
	Sunset        time.Time `json:"sunset"`
	WeatherCode   int       `json:"weather_code"`
	Condition     string    `json:"condition"`
}

// HourlyForecast represents a single hour's forecast.
type HourlyForecast struct {
	Time          time.Time `json:"time"`
	Temperature   float64   `json:"temperature"`
	Apparent      float64   `json:"apparent"`
	Humidity      int       `json:"humidity"`
	Precipitation float64   `json:"precipitation"`
	Rain          float64   `json:"rain"`
	Snowfall      float64   `json:"snowfall"`
	WindSpeed     float64   `json:"wind_speed"`
	WindDirection int       `json:"wind_direction"`
	Pressure      float64   `json:"pressure"`
	CloudCover    int       `json:"cloud_cover"`
	Visibility    float64   `json:"visibility"`
	UVIndex       float64   `json:"uv_index"`
	PrecipProb    int       `json:"precip_prob"`
	WeatherCode   int       `json:"weather_code"`
	Condition     string    `json:"condition"`
}

// Forecast represents weather forecast data.
type Forecast struct {
	Location Location         `json:"location"`
	Daily    []DailyForecast  `json:"daily,omitempty"`
	Hourly   []HourlyForecast `json:"hourly,omitempty"`
}

// WeatherCode maps WMO weather codes to human-readable conditions.
var WeatherCode = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Light freezing drizzle",
	57: "Dense freezing drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Light freezing rain",
	67: "Heavy freezing rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// GetCondition returns human-readable condition from weather code.
func GetCondition(code int) string {
	if cond, ok := WeatherCode[code]; ok {
		return cond
	}
	return "Unknown"
}

// WindDirection returns compass direction from degrees.
func WindDirection(degrees int) string {
	dirs := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	// Normalize degrees to 0-360 range
	normalized := degrees % 360
	if normalized < 0 {
		normalized += 360
	}
	idx := int((float64(normalized) + 11.25) / 22.5)
	return dirs[idx%16]
}
