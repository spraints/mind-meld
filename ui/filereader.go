package ui

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spraints/mind-meld/lmsp/lmspsimple"
)

func openfile(path string) (tea.Model, tea.Cmd) {
	f := fileReader{path: path}
	return f, f.read
}

type fileReader struct {
	path      string
	readError error
	data      *lmspsimple.File
}

func (f fileReader) read() tea.Msg {
	programs, err := lmspsimple.Read(f.path)
	if err != nil {
		return err
	}
	return programs
}

func (f fileReader) Init() tea.Cmd {
	return f.read
}

func (f fileReader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		f.readError = msg
	case *lmspsimple.File:
		f.data = msg
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
	start := escape + f.path + "\n"
	if f.readError != nil {
		return start + f.readError.Error() + "\n"
	}
	if f.data == nil {
		return start + loading
	}
	return start + f.data.JSON + "\n"
}
