package mdtt

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	ta.Focus()
	ta.SetValue(value)
	ta.Prompt = ""
	ta.Cursor.Style = cellCursorStyle
	return Cell{textInput: ta, err: nil}
}

func (m Cell) Init() tea.Cmd {
	return textinput.Blink
}

func (m Cell) Update(msg tea.Msg) (Cell, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
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
