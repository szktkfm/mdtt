package main

import (
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/szktkfm/mdtt"
)

var (
	Version = "0.0.1"
	Commit  = ""
	Date    = ""
	BuiltBy = ""
	rootCmd = &cobra.Command{
		Use:     "mdtt [file]",
		Short:   "Markdown Table Editor with TUI",
		Version: "",
		Run: func(cmd *cobra.Command, args []string) {
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				log.Fatal("Cannot parse debug flag", err)
			}

			logger := createLogger(debug)
			if logger != nil {
				defer logger.Close()
			}

			inplace, _ := cmd.Flags().GetBool("inplace")
			if inplace && len(args) == 0 {
				log.Fatal("no input files")
			}
			var model = createModel(args, inplace)

			p := tea.NewProgram(
				model,
				tea.WithoutSignalHandler(),
				tea.WithOutput(
					termenv.NewOutput(os.Stderr),
				),
			)
			if _, err := p.Run(); err != nil {
				log.Fatal("Failed starting the TUI", err)
			}
		},
	}
)

func createModel(args []string, inplace bool) mdtt.Model {

	if !isatty.IsTerminal(os.Stdin.Fd()) {

		content, _ := io.ReadAll(os.Stdin)
		model, err := mdtt.NewUI(
			mdtt.WithMarkdown(content),
		)
		if err != nil {
			log.Fatal(err)
		}
		return model

	} else if len(args) == 0 {

		model, err := mdtt.NewUI()
		if err != nil {
			log.Fatal(err)
		}
		return model

	} else {
		f, _ := os.Open(args[0])
		defer f.Close()
		content, _ := io.ReadAll(f)
		model, err := mdtt.NewUI(
			mdtt.WithMarkdown(content),
			mdtt.WithInplace(inplace),
			mdtt.WithFilePath(args[0]),
		)
		if err != nil {
			log.Fatal(err)
		}
		return model
	}
}

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
		"inplace",
		"i",
		false,
		"in-place update",
	)
	rootCmd.Flags().BoolP(
		"help",
		"h",
		false,
		"help for mdtt",
	)
	lipgloss.SetColorProfile(termenv.ANSI256)
}

func createLogger(debug bool) *os.File {
	if debug {
		newConfigFile, fileErr := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if fileErr == nil {
			log.SetOutput(newConfigFile)
			log.SetTimeFormat(time.Kitchen)
			log.SetReportCaller(true)
			log.SetLevel(log.DebugLevel)
			log.Debug("Logging to debug.log")
		} else {
			log.Print("Failed setting up logging", fileErr)
		}
		return newConfigFile
	}

	return nil
}

func main() {
	Execute()
}
