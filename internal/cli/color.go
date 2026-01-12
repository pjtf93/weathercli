package cli

import (
	"os"

	"github.com/fatih/color"
)

// Color wraps color output helpers.
type Color struct {
	enabled bool
	bold    *color.Color
	cyan    *color.Color
	green   *color.Color
	yellow  *color.Color
	red     *color.Color
	blue    *color.Color
	magenta *color.Color
}

// NewColor creates a color helper.
func NewColor(enabled bool) Color {
	return Color{
		enabled: enabled,
		bold:    color.New(color.Bold),
		cyan:    color.New(color.FgCyan),
		green:   color.New(color.FgGreen),
		yellow:  color.New(color.FgYellow),
		red:     color.New(color.FgRed),
		blue:    color.New(color.FgBlue),
		magenta: color.New(color.FgMagenta),
	}
}

// Bold returns bold text.
func (c Color) Bold(s string) string {
	if !c.enabled {
		return s
	}
	return c.bold.Sprint(s)
}

// Cyan returns cyan text.
func (c Color) Cyan(s string) string {
	if !c.enabled {
		return s
	}
	return c.cyan.Sprint(s)
}

// Green returns green text.
func (c Color) Green(s string) string {
	if !c.enabled {
		return s
	}
	return c.green.Sprint(s)
}

// Yellow returns yellow text.
func (c Color) Yellow(s string) string {
	if !c.enabled {
		return s
	}
	return c.yellow.Sprint(s)
}

// Red returns red text.
func (c Color) Red(s string) string {
	if !c.enabled {
		return s
	}
	return c.red.Sprint(s)
}

// Blue returns blue text.
func (c Color) Blue(s string) string {
	if !c.enabled {
		return s
	}
	return c.blue.Sprint(s)
}

// Magenta returns magenta text.
func (c Color) Magenta(s string) string {
	if !c.enabled {
		return s
	}
	return c.magenta.Sprint(s)
}

// colorEnabled checks if color output should be enabled.
func colorEnabled(noColor bool) bool {
	if noColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return true
}
