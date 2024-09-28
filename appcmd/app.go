package appcmd

type App interface {
	FullName() string
	ProjectDirs() []string
}
