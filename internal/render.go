package internal

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type delegate struct{}

func (d delegate) Height() int  { return 1 }
func (d delegate) Spacing() int { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	task, ok := item.(*Task)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderTask(task, index == m.Index()))
}

var (
	activeStyle  = lipgloss.NewStyle().Foreground(SecondaryColor).Width(60).PaddingRight(2)
	defaultStyle = lipgloss.NewStyle().Foreground(SecondaryGrayColor).Width(60).PaddingRight(2)
	cursorStyle  = lipgloss.NewStyle().Foreground(PrimaryColor)
	alignStyle   = lipgloss.NewStyle().PaddingLeft(1)
)

func renderTask(task *Task, selected bool) string {
	cursor := " "
	style := defaultStyle
	if selected {
		cursor = ">"
		style = activeStyle
	}
	description := style.Render(task.Description)
	symbol, style := ToSymbol(task.State, false)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		cursorStyle.Render(cursor),
		style.Render(symbol),
		alignStyle.Render(description),
	)
}

func renderFocusedTask(task *Task) string {
	symbol, symbolStyle := ToSymbol(task.State, true)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		symbolStyle.Render(symbol),
		alignStyle.Render(activeStyle.Render(task.Description)),
	)
}

var (
	activeHeaderStyle  = lipgloss.NewStyle().Foreground(SecondaryColor).PaddingRight(2)
	defaultHeaderStyle = lipgloss.NewStyle().Foreground(SecondaryGrayColor).PaddingRight(2)
)

func renderHeader(showingArchived bool) string {
	if showingArchived {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			defaultHeaderStyle.Render("| Tasks |"),
			activeHeaderStyle.Render("| Archived |"),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		activeHeaderStyle.Render("| Tasks |"),
		defaultHeaderStyle.Render("| Archived |"),
	)
}
