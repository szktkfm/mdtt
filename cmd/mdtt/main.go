package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/szktkfm/mdtt"
)

func main() {
	m := mdtt.NewRoot()
	if _, err := tea.NewProgram(m, tea.WithoutSignalHandler()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
