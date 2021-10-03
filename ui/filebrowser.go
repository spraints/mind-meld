package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type fileBrowser struct {
	dir     string
	pos     int
	entries []os.DirEntry
	readErr error
}

func (f fileBrowser) Init() tea.Cmd {
	return f.read
}

func (f fileBrowser) read() tea.Msg {
	entries, err := os.ReadDir(f.dir)
	if err != nil {
		return err
	} else {
		return entries
	}
}

func (f fileBrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []os.DirEntry:
		f.entries = msg
		return f, nil
	case error:
		f.readErr = msg
		return f, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "up":
			if f.pos > 0 {
				f.pos--
			}
			return f, nil
		case "down":
			if f.pos < len(f.entries)-1 {
				f.pos++
			}
			return f, nil
		case "enter", "right":
			newf := fileBrowser{
				dir: filepath.Join(f.dir, f.entries[f.pos].Name()),
			}
			return newf, newf.read
		case "ctrl+h", "left", "b":
			newf := fileBrowser{
				dir: filepath.Join(f.dir, ".."),
			}
			return newf, newf.read
		}
	}
	return f, nil
}

func (f fileBrowser) View() string {
	const (
		escape  = "(press ctrl+c to exit, <- to go up)\n"
		errrr   = "!!! ERROR !!!\n"
		loading = "... loading ...\n"
	)

	start := escape + f.dir + "\n"

	switch {
	case f.readErr != nil:
		return start + errrr + f.readErr.Error() + "\n"
	case len(f.entries) == 0:
		return start + loading
	default:
		return start + f.list()
	}
}

func (f fileBrowser) list() string {
	const (
		selected   = "> %s >\n"
		unselected = "  %s\n"
	)

	startPos := f.pos - 2
	nEntries := len(f.entries)
	if nEntries < startPos+2 {
		startPos = nEntries - 5
	}
	if startPos < 0 {
		startPos = 0
	}
	lines := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		pos := i + startPos
		if pos < nEntries {
			if f.pos == pos {
				lines = append(lines, fmt.Sprintf(selected, f.entries[pos].Name()))
			} else {
				lines = append(lines, fmt.Sprintf(unselected, f.entries[pos].Name()))
			}
		}
	}
	return strings.Join(lines, "")
}
