package mdtt

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type (
	errMsg error
)

type WidthMsg struct {
	width int
}

type Cell struct {
	textInput textinput.Model
	err       error
}

func NewCell(value string) Cell {
	ta := textinput.New()
	ta.Placeholder = ""
	// ta.ShowLineNumbers = false
	// ta.SetHeight(2)
	ta.Focus()
	// ta.CharLimit = 156
	ta.SetValue(value)
	ta.Prompt = ""
	ta.Cursor.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	return Cell{textInput: ta, err: nil}
}

func (m Cell) Init() tea.Cmd {
	return textinput.Blink
}

func (m Cell) Update(msg tea.Msg) (Cell, tea.Cmd) {
	var cmd tea.Cmd

	// h := m.textInput.LineInfo().Height
	// fmt.Println(h)

	switch msg := msg.(type) {
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	// h := strings.Count(m.textInput.Value(), "\n") + 2
	// h := m.textInput.LineCount() + 1
	// m.textInput.SetHeight(h)
	width := runewidth.StringWidth(m.textInput.Value()) + 2
	return m, tea.Batch(cmd, updateWidthCmd(width))
}

func updateWidthCmd(width int) tea.Cmd {
	return func() tea.Msg {
		return WidthMsg{width}
	}
}

func (m Cell) View() string {
	return m.textInput.View()
}

func (m Cell) Value() string {
	return m.textInput.Value()
}
