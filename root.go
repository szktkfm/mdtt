package mdtt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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
		log.Debug("msg", "msg", msg)
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
	if file == "" {
		columns := []Column{
			{Title: NewCell(""), Width: 4},
			{Title: NewCell(""), Width: 4},
		}

		rows := []NaiveRow{
			{"", ""},
		}

		t := New(
			WithColumns(columns),
			WithNaiveRows(rows),
			WithFocused(true),
			WithHeight(len(rows)+1),
		)

		s := DefaultStyles()
		t.SetStyles(s)
		return Model{t}
	}

	return parse(file)
}
