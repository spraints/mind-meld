package appcmd

type App interface {
	FullName() string
	RemoteName() string
	ProjectDirs() []string
}
