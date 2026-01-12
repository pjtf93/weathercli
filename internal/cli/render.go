package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pjtf93/weathercli"
)

// RenderCurrentWeather outputs current weather in human or JSON format.
func (a *App) RenderCurrentWeather(w *weathercli.CurrentWeather) error {
	if a.json {
		return json.NewEncoder(a.out).Encode(w)
	}

	loc := w.Location
	locStr := loc.Name
	if loc.Admin1 != "" {
		locStr += ", " + loc.Admin1
	}
	if loc.Country != "" {
		locStr += ", " + loc.Country
	}

	fmt.Fprintf(a.out, "%s\n", a.color.Bold(locStr))
	fmt.Fprintf(a.out, "%s\n\n", a.color.Cyan(w.Time.Format("Mon Jan 2, 2006 15:04 MST")))

	fmt.Fprintf(a.out, "%s %s\n", a.color.Bold("Condition:"), w.Condition)
	fmt.Fprintf(a.out, "%s %s\n", a.color.Bold("Temperature:"), formatTemp(w.Temperature, a.color))
	fmt.Fprintf(a.out, "%s %s\n", a.color.Bold("Feels like:"), formatTemp(w.Apparent, a.color))
	fmt.Fprintf(a.out, "%s %d%%\n", a.color.Bold("Humidity:"), w.Humidity)
	fmt.Fprintf(a.out, "%s %.1f km/h %s\n", a.color.Bold("Wind:"), w.WindSpeed, weathercli.WindDirection(w.WindDirection))
	fmt.Fprintf(a.out, "%s %.0f hPa\n", a.color.Bold("Pressure:"), w.Pressure)
	fmt.Fprintf(a.out, "%s %d%%\n", a.color.Bold("Cloud cover:"), w.CloudCover)

	if w.Precipitation > 0 {
		fmt.Fprintf(a.out, "%s %.1f mm\n", a.color.Bold("Precipitation:"), w.Precipitation)
	}
	if w.Rain > 0 {
		fmt.Fprintf(a.out, "%s %.1f mm\n", a.color.Bold("Rain:"), w.Rain)
	}
	if w.Snowfall > 0 {
		fmt.Fprintf(a.out, "%s %.1f cm\n", a.color.Bold("Snowfall:"), w.Snowfall)
	}
	if w.UVIndex > 0 {
		fmt.Fprintf(a.out, "%s %.1f %s\n", a.color.Bold("UV Index:"), w.UVIndex, formatUVLevel(w.UVIndex))
	}

	return nil
}

// RenderForecast outputs forecast in human or JSON format.
func (a *App) RenderForecast(f *weathercli.Forecast) error {
	if a.json {
		return json.NewEncoder(a.out).Encode(f)
	}

	loc := f.Location
	locStr := loc.Name
	if loc.Admin1 != "" {
		locStr += ", " + loc.Admin1
	}
	if loc.Country != "" {
		locStr += ", " + loc.Country
	}

	fmt.Fprintf(a.out, "%s\n\n", a.color.Bold(locStr))

	if len(f.Daily) > 0 {
		a.renderDailyForecast(f.Daily)
	}

	if len(f.Hourly) > 0 {
		a.renderHourlyForecast(f.Hourly)
	}

	return nil
}

func (a *App) renderDailyForecast(days []weathercli.DailyForecast) {
	for _, day := range days {
		date := day.Date.Format("Mon Jan 2")
		fmt.Fprintf(a.out, "%s\n", a.color.Bold(date))
		fmt.Fprintf(a.out, "  %s %s\n", a.color.Cyan("Condition:"), day.Condition)
		fmt.Fprintf(a.out, "  %s %s (high) / %s (low)\n",
			a.color.Cyan("Temperature:"),
			formatTemp(day.TempMax, a.color),
			formatTemp(day.TempMin, a.color))

		if day.PrecipProb > 0 {
			fmt.Fprintf(a.out, "  %s %d%%\n", a.color.Cyan("Precipitation:"), day.PrecipProb)
		}
		if day.Rain > 0 {
			fmt.Fprintf(a.out, "  %s %.1f mm\n", a.color.Cyan("Rain:"), day.Rain)
		}
		if day.Snowfall > 0 {
			fmt.Fprintf(a.out, "  %s %.1f cm\n", a.color.Cyan("Snowfall:"), day.Snowfall)
		}

		fmt.Fprintf(a.out, "  %s %.1f km/h %s\n",
			a.color.Cyan("Wind:"),
			day.WindSpeedMax,
			weathercli.WindDirection(day.WindDirection))

		sunrise := day.Sunrise.Format("15:04")
		sunset := day.Sunset.Format("15:04")
		fmt.Fprintf(a.out, "  %s %s → %s\n", a.color.Cyan("Sun:"), sunrise, sunset)

		if day.UVIndexMax > 0 {
			fmt.Fprintf(a.out, "  %s %.1f %s\n",
				a.color.Cyan("UV Index:"),
				day.UVIndexMax,
				formatUVLevel(day.UVIndexMax))
		}
		fmt.Fprintln(a.out)
	}
}

func (a *App) renderHourlyForecast(hours []weathercli.HourlyForecast) {
	for _, hour := range hours {
		timeStr := hour.Time.Format("Mon Jan 2 15:04")
		fmt.Fprintf(a.out, "%s\n", a.color.Bold(timeStr))
		fmt.Fprintf(a.out, "  %s %s\n", a.color.Cyan("Condition:"), hour.Condition)
		fmt.Fprintf(a.out, "  %s %s (feels %s)\n",
			a.color.Cyan("Temperature:"),
			formatTemp(hour.Temperature, a.color),
			formatTemp(hour.Apparent, a.color))
		fmt.Fprintf(a.out, "  %s %d%%\n", a.color.Cyan("Humidity:"), hour.Humidity)

		if hour.PrecipProb > 0 {
			fmt.Fprintf(a.out, "  %s %d%%\n", a.color.Cyan("Precipitation chance:"), hour.PrecipProb)
		}
		if hour.Precipitation > 0 {
			fmt.Fprintf(a.out, "  %s %.1f mm\n", a.color.Cyan("Precipitation:"), hour.Precipitation)
		}

		fmt.Fprintf(a.out, "  %s %.1f km/h %s\n",
			a.color.Cyan("Wind:"),
			hour.WindSpeed,
			weathercli.WindDirection(hour.WindDirection))
		fmt.Fprintf(a.out, "  %s %d%%\n", a.color.Cyan("Cloud cover:"), hour.CloudCover)
		fmt.Fprintln(a.out)
	}
}

// RenderLocations outputs location search results.
func (a *App) RenderLocations(locations []weathercli.Location) error {
	if a.json {
		return json.NewEncoder(a.out).Encode(locations)
	}

	for i, loc := range locations {
		parts := []string{loc.Name}
		if loc.Admin1 != "" {
			parts = append(parts, loc.Admin1)
		}
		if loc.Country != "" {
			parts = append(parts, loc.Country)
		}

		fmt.Fprintf(a.out, "%d. %s\n", i+1, a.color.Bold(strings.Join(parts, ", ")))
		fmt.Fprintf(a.out, "   %s %.4f, %.4f\n",
			a.color.Cyan("Coordinates:"),
			loc.Latitude,
			loc.Longitude)

		if loc.Timezone != "" {
			fmt.Fprintf(a.out, "   %s %s\n", a.color.Cyan("Timezone:"), loc.Timezone)
		}
		fmt.Fprintln(a.out)
	}

	return nil
}

// formatTemp colors temperature based on value.
func formatTemp(temp float64, c Color) string {
	str := fmt.Sprintf("%.1f°C", temp)
	switch {
	case temp >= 30:
		return c.Red(str)
	case temp >= 20:
		return c.Yellow(str)
	case temp >= 10:
		return c.Green(str)
	case temp >= 0:
		return c.Cyan(str)
	default:
		return c.Blue(str)
	}
}

// formatUVLevel returns UV index level description.
func formatUVLevel(uv float64) string {
	switch {
	case uv < 3:
		return "(low)"
	case uv < 6:
		return "(moderate)"
	case uv < 8:
		return "(high)"
	case uv < 11:
		return "(very high)"
	default:
		return "(extreme)"
	}
}

// renderVerbose outputs verbose messages.
func (a *App) renderVerbose(format string, args ...interface{}) {
	fmt.Fprintf(a.err, a.color.Cyan("→ ")+format+"\n", args...)
}
