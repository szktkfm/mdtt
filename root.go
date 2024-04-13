package mdtt

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	defaultHeight = 2
)

type Model struct {
	preview bool
	tables  []TableModel
	table   TableModel
	choose  int
	list    headerList
	fpath   string
	inplace bool
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case selectMsg:
		m.preview = false
		m.choose = msg.idx
		m.table = m.tables[msg.idx]
	case quitMsg:
		if !m.preview {
			Write(m)
		}
		return m, tea.Quit
	}

	if m.preview {
		m.list, cmd = m.list.update(msg)
	} else {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.preview {
		return m.list.view()
	} else {
		return m.table.View() + "\n"
	}
}

func NewRoot(opts ...func(*Model)) Model {
	t := NewTableModel(
		WithColumns(DefaultColumns()),
		WithNaiveRows(DefaultRows()),
		WithFocused(true),
		WithHeight(defaultHeight),
		WithStyles(defaultStyles()),
	)
	m := Model{table: t}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func WithMarkdown(b []byte) func(*Model) {
	return func(m *Model) {

		tables := parse(b)
		list := newHeaderList(
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
	}
}

func WithFilePath(f string) func(*Model) {
	return func(m *Model) {
		m.fpath = f
	}
}

func WithInplace(i bool) func(*Model) {
	return func(m *Model) {
		m.inplace = i
	}
}

func DefaultRows() []naiveRow {
	return []naiveRow{
		{"", ""},
	}
}

func DefaultColumns() []column {
	return []column{
		{title: NewCell(""), width: 4},
		{title: NewCell(""), width: 4},
	}
}
