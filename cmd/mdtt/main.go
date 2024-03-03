package main

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/szktkfm/mdtt"
)

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
	BuiltBy = ""
	rootCmd = &cobra.Command{
		Use:     "gh dash",
		Short:   "A gh extension that shows a configurable dashboard of pull requests and issues.",
		Version: "",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().Bool(
		"debug",
		false,
		"passing this flag will allow writing debug output to debug.log",
	)

	rootCmd.Flags().BoolP(
		"help",
		"h",
		false,
		"help for gh-dash",
	)

	rootCmd.Run = func(_ *cobra.Command, _ []string) {
		debug, err := rootCmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatal("Cannot parse debug flag", err)
		}

		model, logger := newFunction(debug)
		if logger != nil {
			defer logger.Close()
		}

		p := tea.NewProgram(
			model,
			tea.WithoutSignalHandler(),
		)
		if _, err := p.Run(); err != nil {
			log.Fatal("Failed starting the TUI", err)
		}
	}
}

func main() {
	Execute()
}

func newFunction(debug bool) (mdtt.Model, *os.File) {
	var loggerFile *os.File

	if debug {
		var fileErr error
		newConfigFile, fileErr := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if fileErr == nil {
			log.SetOutput(newConfigFile)
			log.SetTimeFormat(time.Kitchen)
			log.SetReportCaller(true)
			log.SetLevel(log.DebugLevel)
			log.Debug("Logging to debug.log")
		} else {
			loggerFile, _ = tea.LogToFile("debug.log", "debug")
			log.Print("Failed setting up logging", fileErr)
		}
	}

	return mdtt.NewRoot(), loggerFile

}
