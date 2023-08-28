package editui

import (
	"fmt"
	"notebook/internal/note"
	"notebook/internal/tui/constants"
	"notebook/internal/tui/noteui"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	edit mode = iota
	create
)

type Model struct {
	input    textinput.Model
	textarea textarea.Model
	help     help.Model
	note     *note.Note
	mode     mode
	rootDir  string
}

func New(rootDir string) Model {
	t := textarea.New()
	in := textinput.New()
	in.Focus()

	h := help.New()
	h.ShowAll = true

	return Model{textarea: t, input: in, help: h, rootDir: rootDir}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case noteui.EditMsg:
		var err error
		m.mode = edit
		m.note, err = note.Load(msg.Filepath)
		if err != nil {
			cmds = append(cmds, func() tea.Msg { return noteui.UpdateTableMsg{} })
			cmds = append(cmds, func() tea.Msg { return noteui.BackMsg{} })
			cmds = append(cmds, constants.ErrMsgCmd(fmt.Errorf("can not open file '%s'", msg.Filepath)))
			break
		}
		m.input.SetValue(m.note.Title)
		m.textarea.SetValue(m.note.Text)
		m.textarea.Focus()
		m.input.Blur()
	case noteui.CreateMsg:
		m.mode = create
		m.input.SetValue("")
		m.textarea.SetValue("")
		m.input.Focus()
		m.textarea.Blur()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Save):
			text := m.textarea.Value()
			title := m.input.Value()
			switch m.mode {
			case edit:
				m.note.Text = text
				m.note.ModifiedAt = time.Now()
			case create:
				m.note = note.New(title, text)
			}
			m.note.Save(m.rootDir)
			cmds = append(cmds, func() tea.Msg { return noteui.UpdateTableMsg{} })
			cmds = append(cmds, func() tea.Msg { return noteui.BackMsg{} })
		case m.input.Focused():
			if key.Matches(msg, Keymap.Tab) || key.Matches(msg, constants.Keymap.Enter) {
				m.textarea.Focus()
				m.input.Blur()
			}
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		case m.textarea.Focused():
			if key.Matches(msg, Keymap.Tab) {
				m.input.Focus()
				m.textarea.Blur()
			}
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.input.View(), m.textarea.View(), m.help.View(Keymap))
}
