package internal

import "github.com/charmbracelet/bubbles/key"

type KeyMapList struct {
	Up                 key.Binding
	Down               key.Binding
	Add                key.Binding
	Focus              key.Binding
	Return             key.Binding
	Copy               key.Binding
	MoveStateForward   key.Binding
	MoveStateBackward  key.Binding
	IncreasePriority   key.Binding
	DecreasePriority   key.Binding
	ArchivedTaskToggle key.Binding
	FilterToggle       key.Binding
	Edit               key.Binding
	Delete             key.Binding
	Undo               key.Binding
	Help               key.Binding
	Quit               key.Binding
}

func DefaultKeyMapList() KeyMapList {
	return KeyMapList{
		Up:                 key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "move up")),
		Down:               key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "move down")),
		Add:                key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add task")),
		Focus:              key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "focus task")),
		Return:             key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "return")),
		Copy:               key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "copy task description")),
		MoveStateForward:   key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "toggle task state forward")),
		MoveStateBackward:  key.NewBinding(key.WithKeys("["), key.WithHelp("[", "toggle task state back")),
		IncreasePriority:   key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "increase priority")),
		DecreasePriority:   key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "decrease priority")),
		ArchivedTaskToggle: key.NewBinding(key.WithKeys("!"), key.WithHelp("!", "archive")),
		FilterToggle:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle filter")),
		Edit:               key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit task")),
		Delete:             key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete task")),
		Undo:               key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo delete")),
		Help:               key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		Quit:               key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+q", "quit")),
	}
}

type KeyMapForm struct {
	Submit key.Binding
	Close  key.Binding
}

func DefaultKeyMapForm() KeyMapForm {
	return KeyMapForm{
		Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		Close:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close")),
	}
}
