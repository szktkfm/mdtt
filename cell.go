package mdtt

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

type (
	errMsg error
)

type widthMsg struct {
	width int
}

type cell struct {
	textInput textinput.Model
	err       error
}

func NewCell(value string) cell {
	ta := textinput.New()
	ta.Placeholder = ""
	ta.Focus()
	ta.SetValue(value)
	ta.Prompt = ""
	ta.Cursor.Style = cellCursorStyle
	return cell{textInput: ta, err: nil}
}

func (m cell) update(msg tea.Msg) (cell, tea.Cmd) {
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
		return widthMsg{width}
	}
}

func (m cell) view() string {
	return m.textInput.View()
}

func (m cell) value() string {
	return m.textInput.Value()
}
