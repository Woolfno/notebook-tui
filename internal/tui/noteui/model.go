package noteui

import (
	"fmt"
	"log"
	"notebook/internal/note"
	"notebook/internal/tui/constants"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Choices int

const (
	EMPTY Choices = iota
	CREATE
	OPEN
	DELETE
	EDIT
	EXIT
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table   table.Model
	help    help.Model
	choice  Choices
	rootDir string
	note    *note.Note
}

func New(rootDir string) Model {
	colums := []table.Column{
		{Title: "#", Width: 2},
		{Title: "Title", Width: 20},
		{Title: "Created At", Width: 20},
		{Title: "Modified At", Width: 20},
	}

	rows := loadNotes(rootDir)

	t := table.New(
		table.WithColumns(colums),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := Model{
		table:   t,
		help:    help.New(),
		choice:  EMPTY,
		rootDir: rootDir}
	m.help.ShowAll = true

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case UpdateTableMsg:
		rows := loadNotes(m.rootDir)
		m.table.SetRows(rows)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Create):
			cmds = append(cmds, createNoteCmd(m.rootDir))
			m.choice = CREATE
		case key.Matches(msg, constants.Keymap.Edit):
			if len(m.table.SelectedRow()) > 0 {
				m.choice = EDIT
				cmds = append(cmds, editNoteCmd(m.rootDir, m.table.SelectedRow()[1]))
			}
		case key.Matches(msg, constants.Keymap.Enter):
			if len(m.table.SelectedRow()) > 0 {
				var err error
				m.choice = OPEN
				path := filepath.Join(m.rootDir, m.table.SelectedRow()[1])
				m.note, err = note.Load(path)
				if err != nil {
					cmds = append(cmds, constants.ErrMsgCmd(fmt.Errorf("can not open file '%s'", path)))
					cmds = append(cmds, func() tea.Msg { return UpdateTableMsg{} })
				}
			}
		case key.Matches(msg, constants.Keymap.Delete):
			if len(m.table.SelectedRow()) > 0 {
				m.choice = DELETE
				cmds = append(cmds, deleteNoteCmd(m.rootDir, m.table.SelectedRow()[1]))
			}
		case key.Matches(msg, constants.Keymap.Quit):
			m.choice = EXIT
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			m.choice = EMPTY
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.choice {
	case OPEN:
		if m.note != nil {
			return baseStyle.Render(fmt.Sprintf("%s\n%s", m.note.Title, m.note.Text))
		}
	case EXIT:
		return "Goodbye!"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		baseStyle.Render(m.table.View()),
		m.help.View(constants.Keymap),
	)
}

func loadNotes(rootDir string) []table.Row {
	const op = "noteui.loadNotes"

	dir, err := os.Open(rootDir)
	if err != nil {
		log.Fatalf("%s: error open dir: %s", op, err)
	}
	defer dir.Close()

	files, _ := dir.ReadDir(0)
	var rows []table.Row
	index := 1
	for _, file := range files {
		if !file.IsDir() {
			n, err := note.Load(filepath.Join(rootDir, file.Name()))
			if err != nil {
				log.Fatalf("%s: error load note: %s", op, err)
			}
			rows = append(rows, []string{
				strconv.Itoa(index),
				n.Title,
				n.CreatedAt.Format(note.TimeFormat),
				n.ModifiedAt.Format(note.TimeFormat),
			})
			index++
		}
	}
	return rows
}

func createNoteCmd(rootDir string) tea.Cmd {
	return func() tea.Msg {
		return CreateMsg{RootDir: rootDir}
	}
}

func deleteNoteCmd(rootDir string, titleNote string) tea.Cmd {
	return func() tea.Msg {
		const op = "noteui.deleteNoteCmd"

		path := filepath.Join(rootDir, titleNote)
		if err := os.Remove(path); err != nil {
			log.Printf("%s: fail remove note: %s", op, err)
		}
		return UpdateTableMsg{}
	}
}

func editNoteCmd(rootDir string, titleNote string) tea.Cmd {
	return func() tea.Msg {
		return EditMsg{Filepath: filepath.Join(rootDir, titleNote)}
	}
}
