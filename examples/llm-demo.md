# LLM Usage Examples

## Quick Reference for AI Agents

### Get Current Weather

```bash
# Command
weathercli current "Tokyo" --json

# Sample Output (parsed)
{
  "temperature": 3.0,
  "condition": "Clear sky",
  "humidity": 34,
  "wind_speed": 7.1
}
```

### Get Forecast

```bash
# Command
weathercli forecast "Paris" --days 3 --json

# Sample Output (first day)
{
  "date": "2024-01-12",
  "temp_max": 12.1,
  "temp_min": 4.3,
  "condition": "Slight rain",
  "precip_prob": 75
}
```

### Search Locations

```bash
# Command
weathercli search "Springfield" --json

# Sample Output (disambiguates multiple matches)
[
  {"name": "Springfield", "country": "United States", "admin1": "Illinois"},
  {"name": "Springfield", "country": "United States", "admin1": "Massachusetts"},
  {"name": "Springfield", "country": "United States", "admin1": "Missouri"}
]
```

## Common LLM Queries

### "What's the weather like in [city]?"

```bash
weathercli current "London" --json | jq '{
  location: .location.name,
  temperature: .temperature,
  condition: .condition,
  feels_like: .apparent
}'
```

**Response format:**
```json
{
  "location": "London",
  "temperature": 10.3,
  "condition": "Overcast",
  "feels_like": 8.1
}
```

### "Will it rain tomorrow in [city]?"

```bash
weathercli forecast "Seattle" --days 2 --json | jq '.daily[1] | {
  date: .date,
  condition: .condition,
  rain_probability: .precip_prob,
  rain_amount: .rain
}'
```

**Response format:**
```json
{
  "date": "2024-01-13",
  "condition": "Moderate rain",
  "rain_probability": 85,
  "rain_amount": 5.2
}
```

### "What's the warmest day this week in [city]?"

```bash
weathercli forecast "Barcelona" --days 7 --json | jq '
  .daily 
  | max_by(.temp_max) 
  | {date: .date, temp_max: .temp_max, condition: .condition}
'
```

**Response format:**
```json
{
  "date": "2024-01-15",
  "temp_max": 18.5,
  "condition": "Mainly clear"
}
```

### "Is it cold in [city] right now?"

```bash
weathercli current "Chicago" --json | jq '{
  temperature: .temperature,
  is_cold: (.temperature < 10),
  is_freezing: (.temperature < 0)
}'
```

**Response format:**
```json
{
  "temperature": -5.3,
  "is_cold": true,
  "is_freezing": true
}
```

### "When is sunrise/sunset in [city]?"

```bash
weathercli forecast "Cairo" --days 1 --json | jq '.daily[0] | {
  sunrise: .sunrise,
  sunset: .sunset
}'
```

**Response format:**
```json
{
  "sunrise": "2024-01-12T06:50:00+02:00",
  "sunset": "2024-01-12T17:20:00+02:00"
}
```

## Integration Patterns

### Function/Tool Definition (OpenAI/Anthropic)

```json
{
  "name": "get_weather",
  "description": "Get current weather conditions for a location",
  "parameters": {
    "type": "object",
    "properties": {
      "location": {
        "type": "string",
        "description": "City name, optionally with state/country (e.g., 'Paris' or 'Paris, France')"
      }
    },
    "required": ["location"]
  }
}
```

### Execution

```python
import json
import subprocess

def get_weather(location: str) -> dict:
    cmd = ["weathercli", "current", location, "--json"]
    result = subprocess.run(cmd, capture_output=True, text=True)
    return json.loads(result.stdout)

def get_forecast(location: str, days: int = 7) -> dict:
    cmd = ["weathercli", "forecast", location, "--days", str(days), "--json"]
    result = subprocess.run(cmd, capture_output=True, text=True)
    return json.loads(result.stdout)
```

### Error Handling

```bash
# Handle location not found
if ! weathercli current "Atlantis" --json 2>/dev/null; then
    echo '{"error": "Location not found"}'
fi

# Disambiguate location
locations=$(weathercli search "Portland" --json | jq -c '.[0:2]')
echo "Did you mean: $locations"
```

## Data Fields Reference

### Current Weather (`current`)
- `temperature`: °C, float
- `apparent`: Feels-like temperature, °C, float
- `humidity`: %, int (0-100)
- `wind_speed`: km/h, float
- `wind_direction`: degrees, int (0-360)
- `condition`: Human-readable string
- `weather_code`: WMO code, int (0-99)
- `precipitation`: mm, float
- `uv_index`: float (0-11+)
- `cloud_cover`: %, int (0-100)

### Daily Forecast (`forecast --days N`)
- `temp_max`, `temp_min`: °C, float
- `apparent_max`, `apparent_min`: Feels-like °C, float
- `precip_prob`: %, int (0-100)
- `rain`, `snowfall`: mm/cm, float
- `sunrise`, `sunset`: ISO8601 timestamp
- `uv_index_max`: float (0-11+)
- `wind_speed_max`: km/h, float

### Hourly Forecast (`forecast --hourly --hours N`)
- Same as current weather plus:
- `precip_prob`: Precipitation probability, %, int

## Tips for LLMs

1. **Always use `--json`** for programmatic parsing
2. **Use `jq`** to extract specific fields (pre-parse before returning to user)
3. **Include location context** in responses (city, country)
4. **Round numbers** for readability (e.g., `jq '.temperature | round'`)
5. **Convert units** if needed (°F: `temp * 9/5 + 32`)
6. **Check `precip_prob`** before mentioning rain
7. **Disambiguate** with `search` command when location is ambiguous
8. **Local timezone** is automatically used for all times
