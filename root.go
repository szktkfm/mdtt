package mdtt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Enum of Mode
const (
	NORMAL = iota
	INSERT
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table TableModel
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewRoot() model {
	columns := []Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 20},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 20},
	}

	rows := []NaiveRow{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
	}

	t := New(
		WithColumns(columns),
		WithNaiveRows(rows),
		WithFocused(true),
		// table.WithHeight(7),
	)

	s := DefaultStyles()
	// s.Header = s.Header.
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("240")).
	// 	BorderBottom(true).
	// 	Bold(false)
	// s.Selected = s.Selected.
	// 	Foreground(lipgloss.Color("16")).
	// 	Background(lipgloss.Color("111")).
	// 	Bold(false)
	t.SetStyles(s)

	m := model{t}
	return m
}
