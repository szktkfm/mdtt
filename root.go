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

var (
	defaultHeight = 2
	baseStyle     = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type Model struct {
	preview bool
	tables  []TableModel
	table   TableModel
	choose  int
	list    ListModel
	fpath   string
	inplace bool
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// insertモードでもないのにqが押されたら終了してしまう
			if !m.preview {
				Write(m)
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

func NewRoot(opts ...func(*Model)) Model {
	t := NewTable(
		WithColumns(DefaultColumns()),
		WithNaiveRows(DefaultRows()),
		WithFocused(true),
		WithHeight(defaultHeight),
		WithStyles(DefaultStyles()),
	)
	m := Model{table: t}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func WithMDFile(fpath string) func(*Model) {
	return func(m *Model) {
		tables := parse(fpath)
		list := NewList(
			WithTables(tables),
		)
		m.table = tables[0]
		m.choose = 0
		m.list = list
		m.tables = tables
		if len(tables) == 1 {
			m.preview = false
		} else {
			m.preview = true
		}
		m.fpath = fpath
	}
}

func WithInplace(i bool) func(*Model) {
	return func(m *Model) {
		m.inplace = i
	}
}

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
