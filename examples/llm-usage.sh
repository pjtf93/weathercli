#!/bin/bash
# Example: How an LLM agent might query weather

# Function to get weather for a location
get_weather() {
    local location="$1"
    weathercli current "$location" --json
}

# Function to get forecast
get_forecast() {
    local location="$1"
    local days="${2:-7}"
    weathercli forecast "$location" --days "$days" --json
}

# Function to search locations (for disambiguation)
search_location() {
    local query="$1"
    weathercli search "$query" --json
}

# Example queries an LLM might make

echo "=== Query: What's the weather in Paris? ==="
get_weather "Paris, France" | jq '{
    location: .location.name,
    temperature: .temperature,
    condition: .condition,
    humidity: .humidity
}'

echo ""
echo "=== Query: Will it rain in London this week? ==="
get_forecast "London, UK" --days 7 | jq '[.daily[] | {
    date: .date,
    condition: .condition,
    rain_probability: .precip_prob,
    rain_amount: .rain
}]'

echo ""
echo "=== Query: Is it cold in New York right now? ==="
get_weather "New York, NY" | jq '{
    temperature: .temperature,
    feels_like: .apparent,
    is_cold: (.temperature < 10)
}'

echo ""
echo "=== Query: What's the warmest day this week in Tokyo? ==="
get_forecast "Tokyo, Japan" --days 7 | jq '[.daily[] | {
    date: .date,
    temp_max: .temp_max
}] | max_by(.temp_max)'

echo ""
echo "=== Query: Find Barcelona (disambiguate) ==="
search_location "Barcelona" | jq '.[:3] | .[] | {
    name: .name,
    country: .country,
    coordinates: "\(.latitude),\(.longitude)"
}'
