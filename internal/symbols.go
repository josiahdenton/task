package internal

import "github.com/charmbracelet/lipgloss"

var (
	readyColor     = SecondaryGrayColor
	focusedColor   = PrimaryColor
	holdColor      = PrimaryGrayColor
	completedColor = AccentColor
	urgentColor    = TertiaryColor
)

var (
	readyStyle     = lipgloss.NewStyle().Foreground(readyColor).PaddingRight(1)
	focusedStyle   = lipgloss.NewStyle().Foreground(focusedColor).PaddingRight(1)
	holdSytle      = lipgloss.NewStyle().Foreground(holdColor).PaddingRight(1)
	completedStyle = lipgloss.NewStyle().Foreground(completedColor).PaddingRight(1)
	urgentStyle    = lipgloss.NewStyle().Foreground(urgentColor).PaddingRight(1).Bold(true)
)

func ToSymbol(state TaskState, extraSpace bool) (string, lipgloss.Style) {
	switch state {
	// FIXME: add space by extracting out const to var
	case Ready:
		return " ", readyStyle
	case Focused:
		return " ", focusedStyle
	case Hold:
		return " ", holdSytle
	case Completed:
		return " ", completedStyle
	case Urgent:
		return " ", urgentStyle
	default:
		return "", lipgloss.Style{}
	}
}
