package mindstormsapp

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
	return "LEGO MINDSTORMS Inventor"
}

func (*App) RemoteName() string {
	return "refs/lego/robot-inventor"
}

func (*App) ProjectDirs() []string {
	home := os.Getenv("HOME")
	return []string{
		filepath.Join(home, "Library/Containers/com.lego.retail.mindstorms.robotinventor/Data/Documents/LEGO MINDSTORMS"),
		filepath.Join(home, "Documents/LEGO MINDSTORMS"),
	}
}
