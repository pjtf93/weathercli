package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/pablotovar/weathercli"
)

// App wires CLI output and API access.
type App struct {
	client  *weathercli.Client
	out     io.Writer
	err     io.Writer
	json    bool
	color   Color
	verbose bool
}

// Run executes the CLI with the provided arguments.
func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	root := Root{}
	exitCode := 0
	parser, err := kong.New(
		&root,
		kong.Name("weathercli"),
		kong.Description("Simple weather CLI for humans and LLMs. Get current weather or forecasts for any location."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}),
		kong.Writers(stdout, stderr),
		kong.Exit(func(code int) {
			exitCode = code
			panic(exitSignal{code: code})
		}),
		kong.Vars{"version": Version},
	)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}

	ctx, exited, err := parseWithExit(parser, args, &exitCode)
	if exited {
		return exitCode
	}
	if err != nil {
		if parseErr, ok := err.(*kong.ParseError); ok {
			_ = parseErr.Context.PrintUsage(true)
			_, _ = fmt.Fprintln(stderr, parseErr.Error())
			return 2
		}
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	if root.Global.JSON {
		// JSON output should never include ANSI escapes.
		root.Global.NoColor = true
	}

	client := weathercli.NewClient(weathercli.Options{
		BaseURL:    root.Global.BaseURL,
		GeoBaseURL: root.Global.GeoBaseURL,
		Timeout:    root.Global.Timeout,
	})

	app := &App{
		client:  client,
		out:     stdout,
		err:     stderr,
		json:    root.Global.JSON,
		color:   NewColor(colorEnabled(root.Global.NoColor)),
		verbose: root.Global.Verbose,
	}

	ctx.Bind(app)
	if err := ctx.Run(); err != nil {
		return handleError(stderr, app.color, err)
	}

	return 0
}

type exitSignal struct {
	code int
}

func parseWithExit(parser *kong.Kong, args []string, exitCode *int) (ctx *kong.Context, exited bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if sig, ok := recovered.(exitSignal); ok {
				*exitCode = sig.code
				exited = true
			} else {
				panic(recovered)
			}
		}
	}()

	ctx, err = parser.Parse(args)
	return ctx, false, err
}

func handleError(w io.Writer, c Color, err error) int {
	if exitErr, ok := err.(ExitCodeError); ok {
		return int(exitErr)
	}
	fmt.Fprintf(w, "%s %v\n", c.Red("Error:"), err)
	return 1
}

// Run for CurrentCmd.
func (c *CurrentCmd) Run(app *App) error {
	if app.verbose {
		app.renderVerbose("Fetching current weather for: %s", c.Location)
	}

	ctx := context.Background()
	weather, err := app.client.Current(ctx, c.Location)
	if err != nil {
		return err
	}

	return app.RenderCurrentWeather(weather)
}

// Run for ForecastCmd.
func (c *ForecastCmd) Run(app *App) error {
	days := c.Days
	if c.Hourly && c.Hours > 0 {
		// Convert hours to days for API
		days = (c.Hours + 23) / 24
		if days > 16 {
			days = 16
		}
	}

	if days < 1 || days > 16 {
		return fmt.Errorf("days must be between 1 and 16")
	}

	if app.verbose {
		if c.Hourly {
			app.renderVerbose("Fetching %d-hour forecast for: %s", c.Hours, c.Location)
		} else {
			app.renderVerbose("Fetching %d-day forecast for: %s", days, c.Location)
		}
	}

	ctx := context.Background()
	forecast, err := app.client.Forecast(ctx, c.Location, days, c.Hourly)
	if err != nil {
		return err
	}

	// Limit hourly results
	if c.Hourly && len(forecast.Hourly) > c.Hours {
		forecast.Hourly = forecast.Hourly[:c.Hours]
	}

	return app.RenderForecast(forecast)
}

// Run for SearchCmd.
func (c *SearchCmd) Run(app *App) error {
	if app.verbose {
		app.renderVerbose("Searching locations: %s", c.Query)
	}

	ctx := context.Background()
	locations, err := app.client.SearchLocation(ctx, c.Query)
	if err != nil {
		return err
	}

	if c.Limit > 0 && len(locations) > c.Limit {
		locations = locations[:c.Limit]
	}

	return app.RenderLocations(locations)
}

// Helper to render JSON output.
func renderJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
