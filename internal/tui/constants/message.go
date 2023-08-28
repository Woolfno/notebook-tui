package constants

import tea "github.com/charmbracelet/bubbletea"

type ErrMsg struct{ Err error }

func ErrMsgCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg{Err: err}
	}
}
