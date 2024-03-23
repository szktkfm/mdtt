package mdtt

import (
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
	preview bool
	tables  []TableModel
	table   TableModel
	choose  int
	list    ListModel
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if !m.preview {
				print(m.table)
			}
			return m, tea.Quit
		}
	case SelectMsg:
		m.preview = false
		m.choose = msg.idx
		m.table = m.tables[msg.idx]
	}

	if m.preview {
		m.list, cmd = m.list.Update(msg)
	} else {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.preview {
		return m.list.View()
	} else {
		return baseStyle.Render(m.table.View()) + "\n"
	}
}

func NewRoot(file string) Model {
	if file == "" {
		t := NewTable(
			WithColumns(DefaultColumns()),
			WithNaiveRows(DefaultRows()),
			WithFocused(true),
			WithHeight(defaultHeight),
			WithStyles(DefaultStyles()),
		)
		return Model{table: t}
	}

	tables := parse(file)
	list := NewList(
		WithTables(tables),
	)
	m := Model{
		table:  tables[0],
		list:   list,
		tables: tables,
	}

	if len(tables) == 1 {
		m.preview = false
	} else {
		m.preview = true
	}
	return m
}

var (
	defaultHeight = 3
)

func DefaultRows() []NaiveRow {
	return []NaiveRow{
		{"", ""},
	}
}

func DefaultColumns() []Column {
	return []Column{
		{Title: NewCell(""), Width: 4},
		{Title: NewCell(""), Width: 4},
	}
}
