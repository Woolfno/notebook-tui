package editui

import (
	"notebook/internal/tui/constants"

	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	Save key.Binding
	Tab  key.Binding
	Quit key.Binding
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save, k.Quit},
	}
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Quit}
}

var Keymap = keymap{
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "tab"),
	),
	Quit: constants.Keymap.Back,
}
