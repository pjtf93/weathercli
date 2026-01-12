package weathercli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL    = "https://api.open-meteo.com/v1"
	defaultGeoBaseURL = "https://geocoding-api.open-meteo.com/v1"
	defaultTimeout    = 10 * time.Second
)

// Client handles weather API requests.
type Client struct {
	baseURL    string
	geoBaseURL string
	httpClient *http.Client
}

// Options for creating a new client.
type Options struct {
	BaseURL    string
	GeoBaseURL string
	Timeout    time.Duration
}

// NewClient creates a new weather client.
func NewClient(opts ...Options) *Client {
	opt := Options{
		BaseURL:    defaultBaseURL,
		GeoBaseURL: defaultGeoBaseURL,
		Timeout:    defaultTimeout,
	}
	if len(opts) > 0 {
		if opts[0].BaseURL != "" {
			opt.BaseURL = opts[0].BaseURL
		}
		if opts[0].GeoBaseURL != "" {
			opt.GeoBaseURL = opts[0].GeoBaseURL
		}
		if opts[0].Timeout > 0 {
			opt.Timeout = opts[0].Timeout
		}
	}

	return &Client{
		baseURL:    opt.BaseURL,
		geoBaseURL: opt.GeoBaseURL,
		httpClient: &http.Client{Timeout: opt.Timeout},
	}
}

// SearchLocation finds locations by name.
func (c *Client) SearchLocation(ctx context.Context, query string) ([]Location, error) {
	u, err := url.Parse(c.geoBaseURL + "/search")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("name", query)
	q.Set("count", "10")
	q.Set("language", "en")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("geocoding API error: %d %s", resp.StatusCode, string(body))
	}

	var result struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Country   string  `json:"country"`
			Admin1    string  `json:"admin1"`
			Timezone  string  `json:"timezone"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("location not found: %s", query)
	}

	locations := make([]Location, len(result.Results))
	for i, r := range result.Results {
		locations[i] = Location{
			Name:      r.Name,
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
			Country:   r.Country,
			Admin1:    r.Admin1,
			Timezone:  r.Timezone,
		}
	}

	return locations, nil
}

// Current fetches current weather for a location.
func (c *Client) Current(ctx context.Context, location string) (*CurrentWeather, error) {
	locations, err := c.SearchLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	loc := locations[0]
	return c.CurrentByCoords(ctx, loc.Latitude, loc.Longitude, &loc)
}

// CurrentByCoords fetches current weather by coordinates.
func (c *Client) CurrentByCoords(ctx context.Context, lat, lon float64, loc *Location) (*CurrentWeather, error) {
	u, err := url.Parse(c.baseURL + "/forecast")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current", "temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,rain,snowfall,weather_code,cloud_cover,pressure_msl,surface_pressure,wind_speed_10m,wind_direction_10m,uv_index")
	q.Set("timezone", "auto")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("weather API error: %d %s", resp.StatusCode, string(body))
	}

	var result struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timezone"`
		Current   struct {
			Time              string  `json:"time"`
			Temperature       float64 `json:"temperature_2m"`
			Apparent          float64 `json:"apparent_temperature"`
			Humidity          int     `json:"relative_humidity_2m"`
			Precipitation     float64 `json:"precipitation"`
			Rain              float64 `json:"rain"`
			Snowfall          float64 `json:"snowfall"`
			WeatherCode       int     `json:"weather_code"`
			CloudCover        int     `json:"cloud_cover"`
			Pressure          float64 `json:"pressure_msl"`
			WindSpeed         float64 `json:"wind_speed_10m"`
			WindDirection     int     `json:"wind_direction_10m"`
			UVIndex           float64 `json:"uv_index"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t, _ := time.Parse(time.RFC3339, result.Current.Time)

	weather := &CurrentWeather{
		Time:          t,
		Temperature:   result.Current.Temperature,
		Apparent:      result.Current.Apparent,
		Humidity:      result.Current.Humidity,
		Precipitation: result.Current.Precipitation,
		Rain:          result.Current.Rain,
		Snowfall:      result.Current.Snowfall,
		WindSpeed:     result.Current.WindSpeed,
		WindDirection: result.Current.WindDirection,
		Pressure:      result.Current.Pressure,
		CloudCover:    result.Current.CloudCover,
		UVIndex:       result.Current.UVIndex,
		WeatherCode:   result.Current.WeatherCode,
		Condition:     GetCondition(result.Current.WeatherCode),
	}

	if loc != nil {
		weather.Location = *loc
		weather.Location.Timezone = result.Timezone
	} else {
		weather.Location = Location{
			Latitude:  result.Latitude,
			Longitude: result.Longitude,
			Timezone:  result.Timezone,
		}
	}

	return weather, nil
}

// Forecast fetches weather forecast for a location.
func (c *Client) Forecast(ctx context.Context, location string, days int, hourly bool) (*Forecast, error) {
	locations, err := c.SearchLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	loc := locations[0]
	return c.ForecastByCoords(ctx, loc.Latitude, loc.Longitude, days, hourly, &loc)
}

// ForecastByCoords fetches forecast by coordinates.
func (c *Client) ForecastByCoords(ctx context.Context, lat, lon float64, days int, hourly bool, loc *Location) (*Forecast, error) {
	u, err := url.Parse(c.baseURL + "/forecast")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("timezone", "auto")
	q.Set("forecast_days", fmt.Sprintf("%d", days))

	if hourly {
		q.Set("hourly", "temperature_2m,relative_humidity_2m,apparent_temperature,precipitation_probability,precipitation,rain,snowfall,weather_code,pressure_msl,cloud_cover,wind_speed_10m,wind_direction_10m,uv_index")
	} else {
		q.Set("daily", "temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,sunrise,sunset,uv_index_max,precipitation_sum,rain_sum,snowfall_sum,precipitation_probability_max,weather_code,wind_speed_10m_max,wind_direction_10m_dominant")
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("weather API error: %d %s", resp.StatusCode, string(body))
	}

	forecast := &Forecast{}

	if loc != nil {
		forecast.Location = *loc
	} else {
		forecast.Location = Location{Latitude: lat, Longitude: lon}
	}

	if hourly {
		var result struct {
			Hourly struct {
				Time          []string  `json:"time"`
				Temperature   []float64 `json:"temperature_2m"`
				Apparent      []float64 `json:"apparent_temperature"`
				Humidity      []int     `json:"relative_humidity_2m"`
				PrecipProb    []int     `json:"precipitation_probability"`
				Precipitation []float64 `json:"precipitation"`
				Rain          []float64 `json:"rain"`
				Snowfall      []float64 `json:"snowfall"`
				WeatherCode   []int     `json:"weather_code"`
				Pressure      []float64 `json:"pressure_msl"`
				CloudCover    []int     `json:"cloud_cover"`
				WindSpeed     []float64 `json:"wind_speed_10m"`
				WindDirection []int     `json:"wind_direction_10m"`
				UVIndex       []float64 `json:"uv_index"`
			} `json:"hourly"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		forecast.Hourly = make([]HourlyForecast, len(result.Hourly.Time))
		for i := range result.Hourly.Time {
			t, _ := time.Parse(time.RFC3339, result.Hourly.Time[i])
			forecast.Hourly[i] = HourlyForecast{
				Time:          t,
				Temperature:   result.Hourly.Temperature[i],
				Apparent:      result.Hourly.Apparent[i],
				Humidity:      result.Hourly.Humidity[i],
				PrecipProb:    result.Hourly.PrecipProb[i],
				Precipitation: result.Hourly.Precipitation[i],
				Rain:          result.Hourly.Rain[i],
				Snowfall:      result.Hourly.Snowfall[i],
				WeatherCode:   result.Hourly.WeatherCode[i],
				Condition:     GetCondition(result.Hourly.WeatherCode[i]),
				Pressure:      result.Hourly.Pressure[i],
				CloudCover:    result.Hourly.CloudCover[i],
				WindSpeed:     result.Hourly.WindSpeed[i],
				WindDirection: result.Hourly.WindDirection[i],
				UVIndex:       result.Hourly.UVIndex[i],
			}
		}
	} else {
		var result struct {
			Daily struct {
				Time          []string  `json:"time"`
				TempMax       []float64 `json:"temperature_2m_max"`
				TempMin       []float64 `json:"temperature_2m_min"`
				ApparentMax   []float64 `json:"apparent_temperature_max"`
				ApparentMin   []float64 `json:"apparent_temperature_min"`
				Sunrise       []string  `json:"sunrise"`
				Sunset        []string  `json:"sunset"`
				UVIndexMax    []float64 `json:"uv_index_max"`
				Precipitation []float64 `json:"precipitation_sum"`
				Rain          []float64 `json:"rain_sum"`
				Snowfall      []float64 `json:"snowfall_sum"`
				PrecipProb    []int     `json:"precipitation_probability_max"`
				WeatherCode   []int     `json:"weather_code"`
				WindSpeedMax  []float64 `json:"wind_speed_10m_max"`
				WindDirection []int     `json:"wind_direction_10m_dominant"`
			} `json:"daily"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		forecast.Daily = make([]DailyForecast, len(result.Daily.Time))
		for i := range result.Daily.Time {
			date, _ := time.Parse("2006-01-02", result.Daily.Time[i])
			sunrise, _ := time.Parse(time.RFC3339, result.Daily.Sunrise[i])
			sunset, _ := time.Parse(time.RFC3339, result.Daily.Sunset[i])

			forecast.Daily[i] = DailyForecast{
				Date:          date,
				TempMax:       result.Daily.TempMax[i],
				TempMin:       result.Daily.TempMin[i],
				ApparentMax:   result.Daily.ApparentMax[i],
				ApparentMin:   result.Daily.ApparentMin[i],
				Sunrise:       sunrise,
				Sunset:        sunset,
				UVIndexMax:    result.Daily.UVIndexMax[i],
				Precipitation: result.Daily.Precipitation[i],
				Rain:          result.Daily.Rain[i],
				Snowfall:      result.Daily.Snowfall[i],
				PrecipProb:    result.Daily.PrecipProb[i],
				WeatherCode:   result.Daily.WeatherCode[i],
				Condition:     GetCondition(result.Daily.WeatherCode[i]),
				WindSpeedMax:  result.Daily.WindSpeedMax[i],
				WindDirection: result.Daily.WindDirection[i],
			}
		}
	}

	return forecast, nil
}
