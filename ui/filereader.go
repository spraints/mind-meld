package ui

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func openfile(path string) (tea.Model, tea.Cmd) {
	f := fileReader{path: path}
	return f, f.read
}

type fileReader struct {
	path string
}

func (f fileReader) read() tea.Msg {
	return nil
}

func (f fileReader) Init() tea.Cmd {
	return f.read
}

func (f fileReader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "ctrl+h", "left":
			return chdir(filepath.Dir(f.path))
		}
	}
	return f, nil
}

func (f fileReader) View() string {
	return escape + f.path + "\n" + "!!! TODO !!!\n"
}
