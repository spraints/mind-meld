package mindstormsapp

type App struct {
}

func New() *App {
	return &App{}
}

func (*App) FullName() string {
	return "LEGO MINDSTORMS Inventor"
}
