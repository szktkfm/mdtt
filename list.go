package mdtt

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerWidth = 30
)

type item string

type selectMsg struct {
	idx int
}

func selectCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return selectMsg{idx: idx}
	}
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := string(i)

	fn := listItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listSelectedItemStyle.Render("> " +
				strings.Replace(strings.Join(s, " "), "\n", "\n  ", -1))
		}
	}

	fmt.Fprint(w, fn(str))
}

type headerList struct {
	list list.Model
}

func (m headerList) update(msg tea.Msg) (headerList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			return m, selectCmd(m.list.Index())
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m headerList) view() string {
	return "\n" + m.list.View()
}

func newHeaderList(opts ...func(*headerList)) headerList {
	listItems := []list.Item{}
	const defaultWidth = 20

	l := list.New(listItems, itemDelegate{}, defaultWidth, 0)
	l.Title = "Select a table:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = listPaginationStyle
	l.Styles.HelpStyle = listHelpStyle

	m := headerList{list: l}

	for _, opt := range opts {
		opt(&m)
	}

	m.list.SetHeight(len(m.list.Items()) * 2)

	return m

}

func WithItems(items []list.Item) func(*headerList) {
	return func(m *headerList) {
		m.list.SetItems(items)
	}
}

func WithTables(tables []TableModel) func(*headerList) {
	return func(m *headerList) {
		var items []list.Item
		for _, t := range tables {
			var headerCells []string
			for _, c := range t.cols {
				headerCells = append(headerCells, strings.Repeat(" ", 4)+c.title.value())
			}
			header := lipgloss.JoinHorizontal(lipgloss.Left, headerCells...)
			header = padOrTruncate(header, headerWidth) + " …"
			items = append(items, item(listHeaderStyle.Render(header)))
		}
		m.list.SetItems(items)
	}
}
