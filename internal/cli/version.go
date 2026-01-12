package cli

import "fmt"

// Version is set by goreleaser at build time.
var Version = "dev"

// VersionFlag is a custom kong flag for --version.
type VersionFlag bool

// BeforeReset is called before the app runs, allowing early exit.
func (v VersionFlag) BeforeReset(app *App, kongCtx interface{}) error {
	if v {
		fmt.Fprintf(app.out, "weathercli %s\n", Version)
		return ExitCodeError(0)
	}
	return nil
}

// ExitCodeError is an error that carries an exit code.
type ExitCodeError int

func (e ExitCodeError) Error() string {
	return fmt.Sprintf("exit %d", e)
}
