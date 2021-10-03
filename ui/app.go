package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return tea.NewProgram(fileBrowser{dir: pwd}).Start()
}
