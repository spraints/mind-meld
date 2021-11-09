package ui

import (
	"bytes"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spraints/mind-meld/lmsdump"
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
	dumpLines []string
	pos       int
}

func (f fileReader) read() tea.Msg {
	programs, err := lmspsimple.Read(f.path)
	if err != nil {
		return err
	}
	//ioutil.WriteFile("project.json", programs.Raw, 0644)

	f.data = programs

	var dumped bytes.Buffer
	if err := lmsdump.Dump(&dumped, programs.Project); err != nil {
		return err
	}
	f.dumpLines = strings.Split(dumped.String(), "\n")

	return f
}

func (f fileReader) Init() tea.Cmd {
	return f.read
}

func (f fileReader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		f.readError = msg
	case fileReader:
		f = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			if f.pos+1 < len(f.dumpLines) {
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
	if f.pos == 0 {
		start += "\n"
	} else {
		start += "--^^ more ^^--\n"
	}
	lines := f.dumpLines[f.pos:]
	if len(lines) > 20 {
		lines = lines[:20]
	}
	start += strings.Join(lines, "\n")
	if f.pos >= len(f.dumpLines)-20 {
		start += "\n\n"
	} else {
		start += "\n--vv more vv--\n"
	}
	return start
}
