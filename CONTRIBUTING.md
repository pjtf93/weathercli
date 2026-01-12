# Contributing to weathercli

Thanks for your interest in contributing! 🎉

## Quick Start

1. **Fork & Clone**
   ```bash
   git clone https://github.com/pablotovar/weathercli.git
   cd weathercli
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Run Tests**
   ```bash
   # Unit tests (fast)
   go test -short ./...
   
   # All tests including integration
   go test ./...
   
   # With coverage
   go test -cover ./...
   ```

4. **Build & Test CLI**
   ```bash
   go build -o weathercli ./cmd/weathercli
   ./weathercli current "London"
   ```

## Development Workflow

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feat/my-feature
   ```

2. Make your changes and add tests

3. Run linter:
   ```bash
   golangci-lint run
   ```

4. Ensure all tests pass:
   ```bash
   go test ./...
   ```

5. Commit with conventional commits:
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve bug"
   git commit -m "docs: update readme"
   ```

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Adding or updating tests
- `refactor:` Code refactoring
- `style:` Code style changes (formatting, etc.)
- `chore:` Maintenance tasks
- `perf:` Performance improvements

### Code Style

- Follow standard Go conventions (`gofmt`, `goimports`)
- Add godoc comments for exported functions/types
- Keep functions focused and testable
- Handle errors explicitly (no silent failures)

### Testing

- Add tests for new features
- Maintain or improve code coverage
- Use `testing.Short()` for fast unit tests
- Integration tests can hit real API (they're skipped with `-short`)

### Pull Requests

1. Update CHANGELOG.md with your changes
2. Ensure CI passes (tests, linting, build)
3. Request review
4. Address feedback

## Project Structure

```
weathercli/
├── cmd/weathercli/     # CLI entry point
├── internal/cli/       # CLI implementation (commands, rendering)
├── client.go           # API client
├── types.go            # Data structures
├── client_test.go      # Tests
└── examples/           # Usage examples
```

## API Integration

We use [Open-Meteo API](https://open-meteo.com/):
- Free, no API key required
- Rate limits: Be respectful, cache when possible
- Timezone handling: Always parse with location timezone

## Questions?

- Open an issue for bugs or feature requests
- Start a discussion for questions or ideas

## License

By contributing, you agree your contributions will be licensed under MIT License.
