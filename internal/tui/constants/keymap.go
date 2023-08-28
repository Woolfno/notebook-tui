package constants

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Create key.Binding
	Edit   key.Binding
	Enter  key.Binding
	Delete key.Binding
	Quit   key.Binding
	Back   key.Binding
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Create, k.Edit, k.Delete},
		{k.Delete, k.Back},
		{k.Quit},
	}
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

var Keymap = keymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
