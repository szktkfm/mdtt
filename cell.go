package mdtt

import (
	"fmt"

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
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetValue(value)
	return Cell{textInput: ti, err: nil}
}

func (m Cell) Init() tea.Cmd {
	return textinput.Blink
}

func (m Cell) Update(msg tea.Msg) (Cell, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	width := runewidth.StringWidth(m.textInput.Value())
	return m, tea.Batch(cmd, updateWidthCmd(width))
}

func updateWidthCmd(width int) tea.Cmd {
	return func() tea.Msg {
		return WidthMsg{width}
	}
}

func (m Cell) View() string {
	return fmt.Sprintf(
		m.textInput.View(),
	) + "\n"
}

func (m Cell) Value() string {
	return m.textInput.Value()
}
