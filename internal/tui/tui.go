package tui

import (
	"log"
	"notebook/internal/tui/constants"
	"notebook/internal/tui/editui"
	"notebook/internal/tui/noteui"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionStata int

const (
	notesView sessionStata = iota
	createView
	editView
)

type MainModel struct {
	notes   tea.Model
	editor  tea.Model
	state   sessionStata
	message string
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""

	switch msg := msg.(type) {
	case noteui.CreateMsg:
		m.state = editView
	case noteui.EditMsg:
		m.state = editView
	case noteui.UpdateTableMsg:
		m.notes, _ = m.notes.Update(msg)
	case noteui.BackMsg:
		m.state = notesView
	case constants.ErrMsg:
		m.message = msg.Err.Error()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Back):
			m.state = notesView
		}
	}

	switch m.state {
	case editView:
		newEdit, newCmd := m.editor.Update(msg)
		editModel, ok := newEdit.(editui.Model)
		if !ok {
			log.Fatal("problem with edit mode")
		}
		m.editor = editModel
		cmd = newCmd
	case notesView:
		newNotes, newCmd := m.notes.Update(msg)
		notesModel, ok := newNotes.(noteui.Model)
		if !ok {
			log.Fatal("problem with notes model")
		}
		m.notes = notesModel
		cmd = newCmd
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var view string
	switch m.state {
	case editView:
		view = m.editor.View()
	default:
		view = m.notes.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, view, m.message)
}

type Tui struct {
	model tea.Model
}

func New(rootDir string) *Tui {
	d, err := os.Open(rootDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(rootDir, os.ModeDir); err != nil {
				log.Fatalf("can not create directory: %s", rootDir)
			}
		} else {
			log.Fatalf("failed open directory: %s", err)
		}
	}
	defer d.Close()

	model := MainModel{
		state:  notesView,
		notes:  noteui.New(rootDir),
		editor: editui.New(rootDir),
	}

	return &Tui{model: model}
}

func (t *Tui) Run() error {
	program := tea.NewProgram(t.model)
	_, err := program.Run()
	return err
}
