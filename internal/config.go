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
	ToggleArchived     key.Binding
	Edit               key.Binding
	Delete             key.Binding
	Undo               key.Binding
	Help               key.Binding
	Quit               key.Binding
}

func DefaultKeyMapList() KeyMapList {
	return KeyMapList{
		Up:                 key.NewBinding(key.WithKeys("k/↑", "up"), key.WithHelp("k", "up")),
		Down:               key.NewBinding(key.WithKeys("j/↓", "down"), key.WithHelp("j", "down")),
		Add:                key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Focus:              key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "focus")),
		Return:             key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "return")),
		Copy:               key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "copy")),
		MoveStateForward:   key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "toggle status")),
		MoveStateBackward:  key.NewBinding(key.WithKeys("["), key.WithHelp("[", "toggle status")),
		IncreasePriority:   key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "priority")),
		DecreasePriority:   key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "priority")),
		ArchivedTaskToggle: key.NewBinding(key.WithKeys("!"), key.WithHelp("!", "archive")),
		ToggleArchived:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle archived")),
		Edit:               key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Delete:             key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Undo:               key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo")),
		Help:               key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:               key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	}
}

func (k KeyMapList) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMapList) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Focus, k.Return, k.ToggleArchived},
		{k.Add, k.Edit, k.Copy, k.Undo, k.ArchivedTaskToggle},
		{k.MoveStateForward, k.MoveStateBackward, k.IncreasePriority, k.DecreasePriority},
		{k.Help, k.Quit},
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
