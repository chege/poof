package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

var (
	NoColor           = false
	Output  io.Writer = os.Stdout
)

func Init(noColor bool) {
	NoColor = noColor
	if noColor || !isatty.IsTerminal(os.Stdout.Fd()) {
		color.NoColor = true
	}
}

func Success(format string, a ...interface{}) {
	fmt.Fprintf(Output, "%s %s\n", color.GreenString("✓"), fmt.Sprintf(format, a...))
}

func Warning(format string, a ...interface{}) {
	fmt.Fprintf(Output, "%s %s\n", color.YellowString("!"), fmt.Sprintf(format, a...))
}

func Error(format string, a ...interface{}) {
	fmt.Fprintf(Output, "%s %s\n", color.RedString("✗"), fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
	fmt.Fprintf(Output, "  %s\n", fmt.Sprintf(format, a...))
}

func Bold(s string) string {
	return color.New(color.Bold).Sprint(s)
}
