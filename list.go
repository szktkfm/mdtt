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
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(lipgloss.Color("170"))
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item string

type SelectMsg struct {
	idx int
}

func SelectCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return SelectMsg{idx: idx}
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

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " +
				strings.Replace(strings.Join(s, " "), "\n", "\n  ", -1))
		}
	}

	fmt.Fprint(w, fn(str))
}

type ListModel struct {
	list list.Model
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			return m, SelectCmd(m.list.Index())
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return m.list.View()
}

func NewList(opts ...func(*ListModel)) ListModel {
	listItems := []list.Item{}
	const defaultWidth = 20

	l := list.New(listItems, itemDelegate{}, defaultWidth, 0)
	l.Title = "Select a table:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := ListModel{list: l}

	for _, opt := range opts {
		opt(&m)
	}

	m.list.SetHeight(len(m.list.Items()) * 2)

	return m

}

func WithItems(items []list.Item) func(*ListModel) {
	return func(m *ListModel) {
		m.list.SetItems(items)
	}
}

var headerStyle = lipgloss.NewStyle().Bold(true).Padding(0, 2).
	Border(lipgloss.NormalBorder(), true, false, true, false).
	BorderForeground(lipgloss.Color("240")).
	Foreground(lipgloss.Color("249"))

func WithTables(tables []TableModel) func(*ListModel) {
	return func(m *ListModel) {
		var items []list.Item
		for _, t := range tables {
			var header []string
			for _, c := range t.cols {
				header = append(header,
					headerStyle.Render(c.Title.Value()))
			}
			header = append(header,
				headerStyle.Render("…"))

			items = append(items, item(lipgloss.JoinHorizontal(lipgloss.Left, header...)))
		}
		m.list.SetItems(items)
	}
}
