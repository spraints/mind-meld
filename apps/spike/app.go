package spike

import (
	"os"
	"path/filepath"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (*App) FullName() string {
	return "SPIKE"
}

func (*App) ProjectDirs() []string {
	home := os.Getenv("HOME")
	return []string{
		filepath.Join(home, "Library/Containers/com.lego.education.spikenext/Data/Documents/LEGO Education SPIKE"),
	}
}
