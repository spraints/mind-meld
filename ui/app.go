package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Opts struct {
	Workdir string
}

func Run(opts Opts) error {
	startdir := opts.Workdir
	if startdir == "" {
		startdir = "."
	}
	return tea.NewProgram(fileBrowser{dir: startdir}).Start()
}
