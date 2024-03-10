package mdtt

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/termenv"
)

// TableModel defines a state for the table widget.
type TableModel struct {
	KeyMap KeyMap

	cols    []Column
	rows    []Row
	cursor  Cursor
	focus   bool
	styles  Styles
	prevKey string

	viewport viewport.Model
	start    int
	end      int
	mode     int

	register Register
}

type Cursor struct {
	x int
	y int
}

// Row represents one line in the table.
type Row []Cell

type NaiveRow []string

// Column defines the table structure.
type Column struct {
	Title Cell
	Width int
}

type Register interface{}

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu.
type KeyMap struct {
	LineUp       key.Binding
	LineDown     key.Binding
	Right        key.Binding
	Left         key.Binding
	AddRow       key.Binding
	DelRow       key.Binding
	Yank         key.Binding
	Paste        key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding
	InsertMode   key.Binding
	NormalMode   key.Binding
}

func init() {
	// DefaultKeyMap = DefaultKeyMap()
	termenv.SetDefaultOutput(termenv.NewOutput(os.Stderr))
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		AddRow: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "Add row"),
		),
		DelRow: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Add row"),
		),
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "Copy"),
		),
		Paste: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Paste"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "½ page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		InsertMode: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert"),
		),
		NormalMode: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl-c", "normal mode"),
		),
	}
}

// Styles contains style definitions for this list component. By default, these
// values are generated by DefaultStyles.
type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
}

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("17")).
			Background(lipgloss.Color("4")),
		Header: lipgloss.NewStyle().Bold(true).Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false),

		Cell: lipgloss.NewStyle().Padding(0, 1),
	}
}

// SetStyles sets the table styles.
func (m *TableModel) SetStyles(s Styles) {
	m.styles = s
	m.UpdateViewport()
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option func(*TableModel)

// New creates a new model for the table widget.
func New(opts ...Option) TableModel {
	m := TableModel{
		cursor:   Cursor{0, 0},
		viewport: viewport.New(0, 10),

		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles(),
		mode:   NORMAL,
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.UpdateViewport()

	return m
}

// WithColumns sets the table columns (headers).
func WithColumns(cols []Column) Option {
	return func(m *TableModel) {
		m.cols = cols
	}
}

// WithRows sets the table rows (data).
func WithRows(rows []Row) Option {
	return func(m *TableModel) {
		m.rows = rows
	}
}

// TODO migration
func WithNaiveRows(rows []NaiveRow) Option {
	return func(m *TableModel) {
		m.rows = make([]Row, len(rows))
		for i, r := range rows {
			m.rows[i] = make(Row, len(r))
			for j, c := range r {
				m.rows[i][j] = NewCell(c)
			}
		}
		m.SetHeight(len(rows))
	}
}

// WithHeight sets the height of the table.
func WithHeight(h int) Option {
	return func(m *TableModel) {
		m.viewport.Height = h
	}
}

// WithWidth sets the width of the table.
func WithWidth(w int) Option {
	return func(m *TableModel) {
		m.viewport.Width = w
	}
}

// WithFocused sets the focus state of the table.
func WithFocused(f bool) Option {
	return func(m *TableModel) {
		m.focus = f
	}
}

// WithStyles sets the table styles.
func WithStyles(s Styles) Option {
	return func(m *TableModel) {
		m.styles = s
	}
}

// WithKeyMap sets the key map.
func WithKeyMap(km KeyMap) Option {
	return func(m *TableModel) {
		m.KeyMap = km
	}
}

// Update is the Bubble Tea update loop.
func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmds []tea.Cmd
	if !m.focus {
		return m, nil
	}

	switch m.mode {
	case NORMAL, HEADER:
		switch msg := msg.(type) {
		case WidthMsg:
			m.UpdateWidth(msg)
		case delPrevKeyMsg:
			m.prevKey = ""
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.KeyMap.LineUp):
				m.MoveUp(1)
			case key.Matches(msg, m.KeyMap.LineDown):
				m.MoveDown(1)
			case key.Matches(msg, m.KeyMap.Right):
				m.MoveRight(1)
			case key.Matches(msg, m.KeyMap.Left):
				m.MoveLeft(1)
			case key.Matches(msg, m.KeyMap.AddRow):
				if m.mode == HEADER {
					return m, nil
				}
				m.AddEmpty()
				m.switchMode(INSERT)
				m.rows[m.cursor.y][m.cursor.x].Update(msg)

			case key.Matches(msg, m.KeyMap.DelRow):
				if m.mode == HEADER {
					return m, nil
				}
				cmd := m.Del()
				cmds = append(cmds, cmd)
			case key.Matches(msg, m.KeyMap.Yank):
				m.Copy()

			case key.Matches(msg, m.KeyMap.Paste):
				m.Paste()

			case key.Matches(msg, m.KeyMap.PageUp):
				m.MoveUp(m.viewport.Height)
			case key.Matches(msg, m.KeyMap.PageDown):
				m.MoveDown(m.viewport.Height)
			case key.Matches(msg, m.KeyMap.HalfPageUp):
				m.MoveUp(m.viewport.Height / 2)
			case key.Matches(msg, m.KeyMap.HalfPageDown):
				m.MoveDown(m.viewport.Height / 2)
			case key.Matches(msg, m.KeyMap.LineDown):
				m.MoveDown(1)
			case key.Matches(msg, m.KeyMap.GotoTop):
				m.GotoTop()
			case key.Matches(msg, m.KeyMap.GotoBottom):
				m.GotoBottom()
			case key.Matches(msg, m.KeyMap.InsertMode):
				if m.mode == HEADER {
					m.switchMode(HEADER_INSERT)
				} else {
					m.switchMode(INSERT)
				}
			}
			m.prevKey = msg.String()
		}
	case INSERT, HEADER_INSERT:
		switch msg := msg.(type) {
		case WidthMsg:
			m.UpdateWidth(msg)
		case delPrevKeyMsg:
			m.prevKey = ""
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.KeyMap.NormalMode):
				if m.mode == HEADER_INSERT {
					m.switchMode(HEADER)
				} else {
					m.switchMode(NORMAL)

				}
			default:
				cmd := m.UpdateFocusedCell(msg)
				cmds = append(cmds, cmd)
			}
			m.prevKey = msg.String()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *TableModel) UpdateFocusedCell(msg tea.KeyMsg) tea.Cmd {
	if m.mode == INSERT {
		newCell, cmd := m.rows[m.cursor.y][m.cursor.x].Update(msg)
		m.rows[m.cursor.y][m.cursor.x] = newCell
		m.UpdateViewport()
		return cmd
	} else if m.mode == HEADER_INSERT {
		newCell, cmd := m.cols[m.cursor.x].Title.Update(msg)
		m.cols[m.cursor.x].Title = newCell
		m.UpdateViewport()
		return cmd
	}
	return nil
}

func (m *TableModel) switchMode(mode int) {
	m.mode = mode
	if mode == INSERT {
		m.styles.Selected = lipgloss.NewStyle().Bold(true)
	} else {
		m.SetStyles(DefaultStyles())
	}
	m.UpdateViewport()
}

func (m TableModel) UpdateWidth(msg WidthMsg) {
	m.cols[m.cursor.x].Width = max(msg.width, m.cols[m.cursor.x].Width)
}

func (m *TableModel) Copy() {
	//TODO
	// Rowのコピーだけを考える。今のところ
	var row Row
	for _, cell := range m.rows[m.cursor.y] {
		row = append(row, NewCell(cell.Value()))
	}
	m.register = row
}

func (m *TableModel) Paste() {
	//TODO
	// Rowのペーストだけを考える。今のところ
	if m.register != nil {
		m.insertRow(m.cursor.y+1, m.register.(Row))
	}
	m.SetHeight(len(m.rows))
	m.MoveDown(1)
	m.UpdateViewport()
}

// Focused returns the focus state of the table.
func (m TableModel) Focused() bool {
	return m.focus
}

// Focus focuses the table, allowing the user to move around the rows and
// interact.
func (m *TableModel) Focus() {
	m.focus = true
	m.UpdateViewport()
}

// Blur blurs the table, preventing selection or movement.
func (m *TableModel) Blur() {
	m.focus = false
	m.UpdateViewport()
}

// View renders the component.
func (m TableModel) View() string {
	return m.headersView() + "\n" + m.viewport.View()
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *TableModel) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.rows))

	// Render only rows from: m.cursor-m.viewport.Height to: m.cursor+m.viewport.Height
	// Constant runtime, independent of number of rows in a table.
	// Limits the number of renderedRows to a maximum of 2*m.viewport.Height
	if m.cursor.y >= 0 {
		m.start = clamp(m.cursor.y-m.viewport.Height, 0, m.cursor.y)
	} else {
		m.start = 0
	}
	m.end = clamp(m.cursor.y+m.viewport.Height, m.cursor.y, len(m.rows))
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m TableModel) SelectedRow() Row {
	if m.cursor.y < 0 || m.cursor.y >= len(m.rows) {
		return nil
	}

	return m.rows[m.cursor.y]
}

// Rows returns the current rows.
func (m TableModel) Rows() []Row {
	return m.rows
}

// SetRows sets a new rows state.
func (m *TableModel) SetRows(r []Row) {
	m.rows = r
}

func (m *TableModel) AddEmpty() {
	if m.prevKey == "v" {
		m.AddColumn()
		return
	}
	newRow := make(Row, len(m.cols))
	for i := range m.cols {
		newRow[i] = NewCell("")
	}
	m.insertRow(m.cursor.y+1, newRow)
	m.SetHeight(len(m.rows))
	m.MoveDown(1)
}

func (m *TableModel) Del() tea.Cmd {
	if m.prevKey == "d" {
		m.deleteRow(m.cursor.y)
		m.SetHeight(len(m.rows))
		m.UpdateViewport()
		return clearPrevKeyCmd()
	} else if m.prevKey == "v" {
		m.deleteColumn(m.cursor.x)
		m.UpdateViewport()
	}

	return clearPrevKeyCmd()
}

func (m *TableModel) AddColumn() {
	var rows []Row
	for i := range m.rows {
		cell := NewCell("")
		rows = append(rows, insertCell(m.rows[i], m.cursor.x+1, cell))
	}
	m.SetRows(rows)

	newCol := insertCol(m.cols, m.cursor.x+1, Column{Title: NewCell(""), Width: 10})
	m.SetColumns(newCol)
	m.MoveRight(1)
}

// SetColumns sets a new columns state.
func (m *TableModel) SetColumns(c []Column) {
	m.cols = c
}

// SetWidth sets the width of the viewport of the table.
func (m *TableModel) SetWidth(w int) {
	m.viewport.Width = w
	m.UpdateViewport()
}

// SetHeight sets the height of the viewport of the table.
func (m *TableModel) SetHeight(h int) {
	// m.viewport.Height = h
	// m.UpdateViewport()
}

// Height returns the viewport height of the table.
func (m TableModel) Height() int {
	return m.viewport.Height
}

// Width returns the viewport width of the table.
func (m TableModel) Width() int {
	return m.viewport.Width
}

// Cursor returns the index of the selected row.
func (m TableModel) Cursor() int {
	return m.cursor.y
}

// SetCursor sets the cursor.y position in the table.
func (m *TableModel) SetCursor(n int) {
	m.cursor.y = clamp(n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// MoveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *TableModel) MoveUp(n int) {
	if m.cursor.y == 0 {
		m.switchMode(HEADER)
	}
	m.cursor.y = clamp(m.cursor.y-n, 0, len(m.rows)-1)

	switch {
	case m.start == 0:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset, 0, m.cursor.y))
	case m.start < m.viewport.Height:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+n, 0, m.cursor.y))
	case m.viewport.YOffset >= 1:
		m.viewport.YOffset = clamp(m.viewport.YOffset+n, 1, m.viewport.Height)
	}
	m.UpdateViewport()
}

// MoveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *TableModel) MoveDown(n int) {
	if m.cursor.y == 0 && m.mode == HEADER {
		m.switchMode(NORMAL)
		return
	}

	m.cursor.y = clamp(m.cursor.y+n, 0, len(m.rows)-1)
	m.UpdateViewport()

	switch {
	case m.end == len(m.rows):
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.viewport.Height))
	case m.cursor.y > (m.end-m.start)/2:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.cursor.y))
	case m.viewport.YOffset > 1:
	case m.cursor.y > m.viewport.YOffset+m.viewport.Height-1:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+1, 0, 1))
	}
}

// MoveRight moves the selection right by any number of rows.
func (m *TableModel) MoveRight(n int) {
	m.cursor.x = clamp(m.cursor.x+n, 0, len(m.cols)-1)
	m.UpdateViewport()
	//TODO viewport
}

// MoveRight moves the selection right by any number of rows.
func (m *TableModel) MoveLeft(n int) {
	m.cursor.x = clamp(m.cursor.x-n, 0, len(m.cols)-1)
	m.UpdateViewport()
	//TODO viewport
}

// GotoTop moves the selection to the first row.
func (m *TableModel) GotoTop() {
	m.MoveUp(m.cursor.y)
}

// GotoBottom moves the selection to the last row.
func (m *TableModel) GotoBottom() {
	m.MoveDown(len(m.rows))
}

//TODO:

// FromValues create the table rows from a simple string. It uses `\n` by
// default for getting all the rows and the given separator for the fields on
// each row.
// func (m *TableModel) FromValues(value, separator string) {
// 	rows := []Row{}
// 	for _, line := range strings.Split(value, "\n") {
// 		r := Row{}
// 		for _, field := range strings.Split(line, separator) {
// 			r = append(r, field)
// 		}
// 		rows = append(rows, r)
// 	}

// 	m.SetRows(rows)
// }

func (m TableModel) headersView() string {
	// selectしたheaderをstyleをinheritしてview
	var s = make([]string, 0, len(m.cols))
	for i, col := range m.cols {
		var style lipgloss.Style
		if i == m.cursor.x && m.mode == HEADER {
			style = m.styles.Selected.PaddingRight(col.Width - len(col.Title.Value()))
		} else {
			style = lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		}

		var renderedCell string
		if i == m.cursor.x && m.mode == HEADER_INSERT {
			renderedCell = style.Render(col.Title.View())
		} else {
			renderedCell = style.Render(col.Title.Value())
		}
		s = append(s, m.styles.Header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *TableModel) renderRow(rowID int) string {
	var s = make([]string, 0, len(m.cols))
	for i, cell := range m.rows[rowID] {
		style := lipgloss.NewStyle().
			Width(m.cols[i].Width).
			MaxWidth(m.cols[i].Width)
			// Inline(true)

		var renderedCell string
		isSelected := i == m.cursor.x &&
			rowID == m.cursor.y &&
			m.mode != HEADER &&
			m.mode != HEADER_INSERT
		isInsertMode := m.mode == INSERT

		var value string
		if isInsertMode && isSelected {
			value = cell.View()
			// log.Debug("インサート", "test", value)
		} else {
			value = cell.Value()
			// log.Debug("", "test", value)
		}

		if isSelected {
			renderedCell = m.styles.Selected.Render(style.Render(value))
		} else {
			renderedCell = m.styles.Cell.Render(style.Render(value))
		}

		s = append(s, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, s...)

	return row
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}

func (m *TableModel) insertRow(idx int, row Row) {
	var rows []Row
	if len(m.rows) == idx {
		rows = append(m.rows, row)
	} else {
		rows = append(m.rows[:idx+1], m.rows[idx:]...)
	}
	rows[idx] = row
	m.SetRows(rows)
}

func (m *TableModel) deleteRow(idx int) {
	var rows []Row
	if len(m.rows) == idx {
		rows = m.rows[:idx-1]
	} else {
		rows = append(m.rows[:idx], m.rows[idx+1:]...)
	}
	m.SetRows(rows)
}

func (m *TableModel) deleteColumn(idx int) {
	var cols []Column
	if len(m.cols) == idx {
		cols = m.cols[:idx-1]
	} else {
		cols = append(m.cols[:idx], m.cols[idx+1:]...)
	}
	m.SetColumns(cols)

	var rows []Row
	for i := range m.rows {
		rows = append(rows, deleteCell(m.rows[i], idx))
	}
	m.SetRows(rows)
}

func insertCell(r Row, idx int, cell Cell) Row {
	var row Row
	if len(r) == idx {
		row = append(r, cell)
	} else {
		row = append(r[:idx+1], r[idx:]...)
	}
	row[idx] = cell
	return row
}

func insertCol(c []Column, idx int, col Column) []Column {
	var newCol []Column
	if len(c) == idx {
		newCol = append(c, col)
	} else {
		newCol = append(c[:idx+1], c[idx:]...)
	}
	newCol[idx] = col
	return newCol
}

func deleteCell(r Row, idx int) Row {
	var row Row
	if len(r) == idx {
		row = r[:idx-1]
	} else {
		row = append(r[:idx], r[idx+1:]...)
	}
	return row
}

type delPrevKeyMsg struct{}

func clearPrevKeyCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return delPrevKeyMsg{}
	})
}
