package mdtt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Enum of Mode
const (
	NORMAL = iota
	INSERT
	HEADER
	HEADER_INSERT
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table TableModel
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			print(m.table)
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func print(m TableModel) {
	var sb strings.Builder
	var width int

	for _, c := range m.cols {
		sb.WriteString("|")
		sb.WriteString(PadOrTruncate(c.Title.Value(), c.Width))
		width += c.Width
	}
	sb.WriteString("|\n")

	for _, c := range m.cols {
		sb.WriteString("|")
		sb.WriteString(strings.Repeat("-", c.Width))
	}
	sb.WriteString("|\n")

	for _, row := range m.rows {
		for i, c := range row {
			sb.WriteString("|")
			sb.WriteString(PadOrTruncate(c.Value(), m.cols[i].Width))
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("\n")

	fmt.Print(sb.String())
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewRoot(file string) Model {
	m := parse(file)

	// columns := []Column{
	// 	{Title: NewCell("Rank"), Width: 4},
	// 	{Title: NewCell("City"), Width: 20},
	// 	{Title: NewCell("Country"), Width: 10},
	// 	{Title: NewCell("Population"), Width: 20},
	// }

	// rows := []NaiveRow{
	// 	{"1", "Tokyo", "Japan", "37,274,000"},
	// 	{"2", "Delhi", "India", "32,065,760"},
	// 	{"3", "Shanghai", "China", "28,516,904"},
	// 	{"4", "Dhaka", "Bangladesh", "22,478,116"},
	// 	{"5", "São Paulo", "Brazil", "22,429,800"},
	// 	{"6", "Mexico City", "Mexico", "22,085,140"},
	// }

	// t := New(
	// 	WithColumns(columns),
	// 	WithNaiveRows(rows),
	// 	WithFocused(true),
	// 	// table.WithHeight(7),
	// )

	// s := DefaultStyles()
	// // s.Header = s.Header.
	// // 	BorderStyle(lipgloss.NormalBorder()).
	// // 	BorderForeground(lipgloss.Color("240")).
	// // 	BorderBottom(true).
	// // 	Bold(false)
	// // s.Selected = s.Selected.
	// // 	Foreground(lipgloss.Color("16")).
	// // 	Background(lipgloss.Color("111")).
	// // 	Bold(false)
	// t.SetStyles(s)

	// m := Model{t}
	return m
}
