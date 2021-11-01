package ui

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
	pos       int
}

func (f fileReader) read() tea.Msg {
	programs, err := lmspsimple.Read(f.path)
	if err != nil {
		return err
	}
	ioutil.WriteFile("project.json", programs.Raw, 0644)
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
		case "down":
			if f.pos+1 < len(f.data.Project.Targets) {
				f.pos++
			}
			return f, nil
		case "up":
			if f.pos > 0 {
				f.pos--
			}
			return f, nil
		case "ctrl+c":
			return f, tea.Quit
		case "ctrl+h", "left":
			return chdir(filepath.Dir(f.path))
		case "enter":
			t := targetRender{file: f, target: f.data.Project.Targets[f.pos]}
			return t, t.index
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
	lines := make([]string, 0, 1+len(f.data.Project.Targets))
	lines = append(lines, start)
	for i, target := range f.data.Project.Targets {
		if i == f.pos {
			lines = append(lines, fmt.Sprintf("> %s >\n", target.Name))
		} else {
			lines = append(lines, fmt.Sprintf("  %s\n", target.Name))
		}
	}
	return strings.Join(lines, "")
}