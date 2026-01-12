# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- [2026-01-12 10:50] Open-source infrastructure: MIT license, GitHub Actions CI/CD, comprehensive test suite (80.4% coverage), CONTRIBUTING.md, golangci-lint config
- [2026-01-12 10:50] Integration tests for current weather, forecasts (daily/hourly), and location search
- [2026-01-12 10:50] CI workflows for testing on multiple Go versions and platforms (Linux/Mac/Windows)
- [2026-01-12 10:50] Automated release workflow for multi-platform binary builds

### Fixed
- [2026-01-12 10:50] WindDirection: handle negative degree values correctly (normalize to 0-360 range)
- [2026-01-12 10:38] Timestamp parsing: corrected time/date display from 0001-01-01 to actual values; fixed sunrise/sunset from 00:00 to correct local times. API returns `2006-01-02T15:04` format without timezone suffix—now parse with location timezone instead of RFC3339.

### Added
- Initial release of weathercli
- Current weather command with temperature, humidity, wind, precipitation
- Daily forecast command (up to 16 days)
- Hourly forecast command (up to 384 hours)
- Location search command with geocoding
- JSON output for LLM-friendly parsing
- Color human-readable output (respects NO_COLOR)
- Free API integration with Open-Meteo (no API key required)
- Global coverage with automatic timezone detection
- Temperature color-coding (hot to cold: red, yellow, green, cyan, blue)
- Weather condition descriptions from WMO codes
- Wind direction compass points
- UV index with severity levels
- Sunrise/sunset times in daily forecasts
- Precipitation probability and amounts
