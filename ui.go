package mdtt

import (
	"fmt"

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

type Option func(*Model) error

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

func NewUI(opts ...Option) (Model, error) {
	t := NewTableModel(
		WithColumns(DefaultColumns()),
		WithNaiveRows(DefaultRows()),
		WithFocused(true),
		WithHeight(defaultHeight),
		WithStyles(defaultStyles()),
	)
	m := Model{table: t}

	for _, opt := range opts {
		err := opt(&m)
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func WithMarkdown(b []byte) Option {
	return func(m *Model) error {
		tables := parse(b)
		if len(tables) == 0 {
			return fmt.Errorf("no tables found")
		}
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
		return nil
	}
}

func WithFilePath(f string) Option {
	return func(m *Model) error {
		m.fpath = f
		return nil
	}
}

func WithInplace(i bool) Option {
	return func(m *Model) error {
		m.inplace = i
		return nil
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
