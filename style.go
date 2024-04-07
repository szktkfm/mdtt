package mdtt

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	// cell styles
	cellCursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	// list styles
	listItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	listSelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(lipgloss.Color("170"))

	listPaginationStyle = list.DefaultStyles().
				PaginationStyle.PaddingLeft(4)

	listHelpStyle = list.DefaultStyles().
			HelpStyle.PaddingLeft(4).PaddingBottom(1)

	listHeaderStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("240")).
			Foreground(lipgloss.Color("249"))

	listHeaderCellStyle = lipgloss.NewStyle().Bold(true).Padding(0, 2).
				Border(lipgloss.NormalBorder(), false, false, false, false).
				BorderForeground(lipgloss.Color("240")).
				Foreground(lipgloss.Color("249"))

	// table styles
	tableSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("17")).
				Background(lipgloss.Color("4"))

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)

	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)
)
