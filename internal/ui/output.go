// Package ui provides terminal user interface utilities including colored output and symbols.
package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

var (
	// NoColor can be set to true to disable all color output.
	NoColor = false
	// Output is the destination for all UI messages (defaults to os.Stdout).
	Output io.Writer = os.Stdout
)

// Init initializes the UI package, detecting TTY status and setting color preferences.
func Init(noColor bool) {
	NoColor = noColor
	if noColor || !isatty.IsTerminal(os.Stdout.Fd()) {
		color.NoColor = true
	}
}

// Success prints a green checkmark followed by a message.
func Success(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(Output, "%s %s\n", color.GreenString("✓"), fmt.Sprintf(format, a...))
}

// Warning prints a yellow exclamation mark followed by a message.
func Warning(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(Output, "%s %s\n", color.YellowString("!"), fmt.Sprintf(format, a...))
}

// Error prints a red cross followed by a message.
func Error(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(Output, "%s %s\n", color.RedString("✗"), fmt.Sprintf(format, a...))
}

// Info prints an indented informational message.
func Info(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(Output, "  %s\n", fmt.Sprintf(format, a...))
}

// Bold returns the string formatted in bold.
func Bold(s string) string {
	return color.New(color.Bold).Sprint(s)
}
