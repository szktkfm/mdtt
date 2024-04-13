package mdtt

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// Enum of Mode
const (
	NORMAL = iota
	INSERT
	HEADER
	HEADER_INSERT
	HELP
)

// TableModel defines a state for the table widget.
type TableModel struct {
	keys     keyMap
	cols     []column
	rows     []row
	cursor   cursor
	focus    bool
	styles   tableStyles
	prevKey  string
	viewport viewport.Model
	start    int
	end      int
	mode     int
	register register
	help     help.Model
}

type cursor struct {
	x int
	y int
}

// row represents one line in the table.
type row []cell

type naiveRow []string

// column defines the table structure.
type column struct {
	title     cell
	width     int
	alignment string
}

type register interface{}

type quitMsg struct{}

func quitCmd() tea.Cmd {
	return func() tea.Msg {
		return quitMsg{}
	}
}

// keyMap defines keybindings. It satisfies to the help.keyMap interface, which
// is used to render the menu.
type keyMap struct {
	lineUp       key.Binding
	lineDown     key.Binding
	right        key.Binding
	left         key.Binding
	addRowCol    key.Binding
	delRowCol    key.Binding
	yank         key.Binding
	paste        key.Binding
	pageUp       key.Binding
	pageDown     key.Binding
	halfPageUp   key.Binding
	halfPageDown key.Binding
	gotoTop      key.Binding
	gotoBottom   key.Binding
	insertMode   key.Binding
	normalMode   key.Binding
	quit         key.Binding
	help         key.Binding
}

// defaultKeyMap returns a default set of keybindings.
func defaultKeyMap() keyMap {
	const spacebar = " "
	return keyMap{
		lineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		lineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		addRowCol: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o/v+o", "add row/column"),
		),
		delRowCol: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("dd/v+d", "add row/column"),
		),
		yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "copy row"),
		),
		paste: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "paste"),
		),
		pageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		pageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		halfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "½ page up"),
		),
		halfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "½ page down"),
		),
		gotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		gotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		insertMode: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert mode"),
		),
		normalMode: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl-c", "normal mode"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.lineUp, k.lineDown, k.left, k.right, k.pageUp, k.pageDown,
			k.halfPageUp, k.halfPageDown, k.gotoTop, k.gotoBottom},
		{k.insertMode, k.normalMode, k.addRowCol, k.delRowCol,
			k.yank, k.paste, k.quit, k.help},
	}
}

// tableStyles contains style definitions for this list component. By default, these
// values are generated by DefaultStyles.
type tableStyles struct {
	header   lipgloss.Style
	cell     lipgloss.Style
	selected lipgloss.Style
}

// defaultStyles returns a set of default style definitions for this table.
func defaultStyles() tableStyles {
	return tableStyles{
		selected: tableSelectedStyle,
		header:   tableHeaderStyle,
		cell:     tableCellStyle,
	}
}

// SetStyles sets the table styles.
func (m *TableModel) SetStyles(s tableStyles) {
	m.styles = s
	m.updateViewport()
}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option func(*TableModel)

// NewTableModel creates a new model for the table widget.
func NewTableModel(opts ...Option) TableModel {
	m := TableModel{
		cursor:   cursor{0, 0},
		viewport: viewport.New(0, 0),

		keys:   defaultKeyMap(),
		styles: defaultStyles(),
		mode:   NORMAL,
		help:   help.New(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.updateViewport()

	return m
}

// WithColumns sets the table columns (headers).
func WithColumns(cols []column) Option {
	return func(m *TableModel) {
		m.cols = cols
	}
}

// WithRows sets the table rows (data).
func WithRows(rows []row) Option {
	return func(m *TableModel) {
		m.rows = rows
	}
}

// TODO migration
func WithNaiveRows(rows []naiveRow) Option {
	return func(m *TableModel) {
		m.rows = make([]row, len(rows))
		for i, r := range rows {
			m.rows[i] = make(row, len(r))
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
func WithStyles(s tableStyles) Option {
	return func(m *TableModel) {
		m.styles = s
	}
}

// WithKeyMap sets the key map.
func WithKeyMap(km keyMap) Option {
	return func(m *TableModel) {
		m.keys = km
	}
}

// Update is the Bubble Tea update loop.
func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.mode {
	case NORMAL, HEADER:
		switch msg := msg.(type) {
		case widthMsg:
			m.updateWidth(msg)
		case delPrevKeyMsg:
			m.prevKey = ""
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.help):
				m.enableAllHelp()
			case key.Matches(msg, m.keys.quit):
				return m, quitCmd()
			case key.Matches(msg, m.keys.lineUp):
				m.moveUp(1)
			case key.Matches(msg, m.keys.lineDown):
				m.moveDown(1)
			case key.Matches(msg, m.keys.right):
				m.moveRight(1)
			case key.Matches(msg, m.keys.left):
				m.moveLeft(1)
			case key.Matches(msg, m.keys.addRowCol):
				m.addEmpty()
				m.switchMode(INSERT)

			case key.Matches(msg, m.keys.delRowCol):
				if m.mode == HEADER {
					return m, nil
				}
				cmd := m.delete()
				cmds = append(cmds, cmd)
			case key.Matches(msg, m.keys.yank):
				m.copy()

			case key.Matches(msg, m.keys.paste):
				m.paste()

			case key.Matches(msg, m.keys.pageUp):
				m.moveUp(m.viewport.Height)
			case key.Matches(msg, m.keys.pageDown):
				m.moveDown(m.viewport.Height)
			case key.Matches(msg, m.keys.halfPageUp):
				m.moveUp(m.viewport.Height / 2)
			case key.Matches(msg, m.keys.halfPageDown):
				m.moveDown(m.viewport.Height / 2)
			case key.Matches(msg, m.keys.lineDown):
				m.moveDown(1)
			case key.Matches(msg, m.keys.gotoTop):
				m.gotoTop()
			case key.Matches(msg, m.keys.gotoBottom):
				m.gotoBottom()
			case key.Matches(msg, m.keys.insertMode):
				if len(m.cols) == 0 {
					return m, nil
				}
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
		case widthMsg:
			m.updateWidth(msg)
		case delPrevKeyMsg:
			m.prevKey = ""
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.normalMode):
				if m.mode == HEADER_INSERT {
					m.switchMode(HEADER)
				} else {
					m.switchMode(NORMAL)
				}
			default:
				cmd := m.updateFocusedCell(msg)
				cmds = append(cmds, cmd)
			}
			m.prevKey = msg.String()
		}
	case HELP:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.help):
				m.disableAllHelp()
			case key.Matches(msg, m.keys.quit):
				return m, quitCmd()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *TableModel) enableAllHelp() {
	m.mode = HELP
	m.help.ShowAll = true
}

func (m *TableModel) disableAllHelp() {
	m.mode = NORMAL
	m.help.ShowAll = false
}

func (m *TableModel) updateFocusedCell(msg tea.KeyMsg) tea.Cmd {
	if m.mode == INSERT {
		newCell, cmd := m.rows[m.cursor.y][m.cursor.x].update(msg)
		m.rows[m.cursor.y][m.cursor.x] = newCell
		m.updateViewport()
		return cmd
	} else if m.mode == HEADER_INSERT {
		newCell, cmd := m.cols[m.cursor.x].title.update(msg)
		m.cols[m.cursor.x].title = newCell
		m.updateViewport()
		return cmd
	}
	return nil
}

func (m *TableModel) switchMode(mode int) {
	m.mode = mode
	if mode == INSERT {
		m.styles.selected = lipgloss.NewStyle().Bold(true)
	} else {
		m.SetStyles(defaultStyles())
	}
	m.updateViewport()
}

func (m TableModel) updateWidth(msg widthMsg) {
	maxWidth := msg.width
	for _, r := range m.rows {
		maxWidth = max(runewidth.StringWidth(r[m.cursor.x].value())+2, maxWidth)
	}
	maxWidth = max(runewidth.StringWidth(m.cols[m.cursor.x].title.value())+2, maxWidth)
	m.cols[m.cursor.x].width = maxWidth
}

func (m *TableModel) copy() {
	if len(m.rows) == 0 {
		return
	}
	var row row
	for _, cell := range m.rows[m.cursor.y] {
		row = append(row, NewCell(cell.value()))
	}
	m.register = row
}

func (m *TableModel) paste() {
	if m.register != nil {
		m.insertRow(m.cursor.y+1, m.register.(row))

		var ro row
		for _, cell := range m.register.(row) {
			ro = append(ro, NewCell(cell.value()))
		}
		m.register = ro
	}
	m.SetHeight(len(m.rows) + 1)
	m.moveDown(1)
	m.updateViewport()
}

// View renders the component.
func (m TableModel) View() string {
	if m.mode == HELP {
		return tableFrameStyle.Render(m.help.View(m.keys))
	}
	return tableFrameStyle.Render(m.headersView()+"\n"+m.viewport.View()) +
		"\n " + m.help.View(m.keys)
}

// updateViewport updates the list content based on the previously defined
// columns and rows.
func (m *TableModel) updateViewport() {
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
func (m TableModel) SelectedRow() row {
	if m.cursor.y < 0 || m.cursor.y >= len(m.rows) {
		return nil
	}

	return m.rows[m.cursor.y]
}

// Rows returns the current rows.
func (m TableModel) Rows() []row {
	return m.rows
}

// SetRows sets a new rows state.
func (m *TableModel) SetRows(r []row) {
	m.rows = r
}

func (m *TableModel) addEmpty() {
	if m.prevKey == "v" {
		m.addColumn()
		return
	}
	newRow := make(row, len(m.cols))
	for i := range m.cols {
		newRow[i] = NewCell("")
	}

	if m.mode == HEADER || len(m.rows) == 0 {
		m.insertRow(0, newRow)
	} else {
		m.insertRow(m.cursor.y+1, newRow)
	}
	m.SetHeight(len(m.rows))
	m.moveDown(1)
}

func (m *TableModel) delete() tea.Cmd {
	if m.prevKey == "d" {
		if len(m.rows) == 0 {
			return nil
		} else if len(m.rows) == 1 {
			m.switchMode(HEADER)
		}
		m.deleteRow(clamp(m.cursor.y, 0, len(m.rows)-1))
		m.cursor.y = clamp(m.cursor.y, 0, len(m.rows)-1)
		m.SetHeight(len(m.rows))
		m.updateViewport()
		return clearPrevKeyCmd()
	} else if m.prevKey == "v" {
		if len(m.cols) == 0 {
			return nil
		}
		m.deleteColumn(m.cursor.x)
		m.cursor.x = clamp(m.cursor.x, 0, len(m.cols)-1)
		m.updateViewport()
	}

	return clearPrevKeyCmd()
}

func (m *TableModel) addColumn() {
	var rows []row
	for i := range m.rows {
		cell := NewCell("")
		rows = append(rows, insertCell(m.rows[i], m.cursor.x+1, cell))
	}
	m.SetRows(rows)

	newCol := insertCol(m.cols, m.cursor.x+1, column{title: NewCell(""), width: 4})
	m.SetColumns(newCol)
	m.moveRight(1)
}

// SetColumns sets a new columns state.
func (m *TableModel) SetColumns(c []column) {
	m.cols = c
}

// SetWidth sets the width of the viewport of the table.
func (m *TableModel) SetWidth(w int) {
	m.viewport.Width = w
	m.updateViewport()
}

// SetHeight sets the height of the viewport of the table.
func (m *TableModel) SetHeight(h int) {
	m.viewport.Height = h
	m.updateViewport()
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
	m.updateViewport()
}

// moveUp moves the selection up by any number of rows.
// It can not go above the first row.
func (m *TableModel) moveUp(n int) {
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
	m.updateViewport()
}

// moveDown moves the selection down by any number of rows.
// It can not go below the last row.
func (m *TableModel) moveDown(n int) {
	if m.cursor.y == 0 && m.mode == HEADER {
		m.switchMode(NORMAL)
		m.updateViewport()
		return
	}

	m.cursor.y = clamp(m.cursor.y+n, 0, len(m.rows)-1)
	m.updateViewport()

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

// moveRight moves the selection right by any number of rows.
func (m *TableModel) moveRight(n int) {
	m.cursor.x = clamp(m.cursor.x+n, 0, len(m.cols)-1)
	m.updateViewport()
}

// MoveRight moves the selection right by any number of rows.
func (m *TableModel) moveLeft(n int) {
	m.cursor.x = clamp(m.cursor.x-n, 0, len(m.cols)-1)
	m.updateViewport()
}

// gotoTop moves the selection to the first row.
func (m *TableModel) gotoTop() {
	m.moveUp(m.cursor.y)
}

// gotoBottom moves the selection to the last row.
func (m *TableModel) gotoBottom() {
	m.moveDown(len(m.rows))
}

func (m TableModel) headersView() string {
	var s = make([]string, 0, len(m.cols))
	for i, col := range m.cols {
		var style lipgloss.Style
		if i == m.cursor.x && m.mode == HEADER {
			style = m.styles.selected.
				Copy().
				PaddingRight(col.width - len(col.title.value()))
		} else {
			style = lipgloss.NewStyle().Width(col.width).MaxWidth(col.width).Inline(true)
		}

		var renderedCell string
		if i == m.cursor.x && m.mode == HEADER_INSERT {
			renderedCell = style.Render(col.title.view())
		} else {
			renderedCell = style.Render(col.title.value())
		}
		s = append(s, m.styles.header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *TableModel) renderRow(rowID int) string {
	var s = make([]string, 0, len(m.cols))
	for i, cell := range m.rows[rowID] {
		style := lipgloss.NewStyle().
			Width(m.cols[i].width).
			MaxWidth(m.cols[i].width)

		var renderedCell string
		isSelected := i == m.cursor.x &&
			rowID == m.cursor.y &&
			m.mode != HEADER &&
			m.mode != HEADER_INSERT
		isInsertMode := m.mode == INSERT

		var value string
		if isInsertMode && isSelected {
			value = cell.view()
		} else {
			value = cell.value()
		}

		if isSelected {
			renderedCell = m.styles.selected.Render(style.Render(value))
		} else {
			renderedCell = m.styles.cell.Render(style.Render(value))
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

func (m *TableModel) insertRow(idx int, ro row) {
	var rows []row
	if len(m.rows) == idx {
		rows = append(m.rows, ro)
	} else {
		rows = append(m.rows[:idx+1], m.rows[idx:]...)
	}
	rows[idx] = ro
	m.SetRows(rows)
}

func (m *TableModel) deleteRow(idx int) {
	var rows []row
	if len(m.rows) == idx {
		rows = m.rows[:idx-1]
	} else {
		rows = append(m.rows[:idx], m.rows[idx+1:]...)
	}
	m.SetRows(rows)
}

func (m *TableModel) deleteColumn(idx int) {
	var cols []column
	if len(m.cols) == idx {
		cols = m.cols[:idx-1]
	} else {
		cols = append(m.cols[:idx], m.cols[idx+1:]...)
	}
	m.SetColumns(cols)

	var rows []row
	for i := range m.rows {
		rows = append(rows, deleteCell(m.rows[i], idx))
	}
	m.SetRows(rows)
}

func insertCell(r row, idx int, cell cell) row {
	var ro row
	if len(r) == idx {
		ro = append(r, cell)
	} else {
		ro = append(r[:idx+1], r[idx:]...)
	}
	ro[idx] = cell
	return ro
}

func insertCol(c []column, idx int, col column) []column {
	var newCol []column
	if len(c) == idx {
		newCol = append(c, col)
	} else {
		newCol = append(c[:idx+1], c[idx:]...)
	}
	newCol[idx] = col
	return newCol
}

func deleteCell(r row, idx int) row {
	var row row
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
