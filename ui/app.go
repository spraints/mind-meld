package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	return tea.NewProgram(newModel()).Start()
}

func newModel() model {
	return model{}
}

type model struct {
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("%#v", msg)
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m model) View() string {
	return "running (press any key to exit)"
}
