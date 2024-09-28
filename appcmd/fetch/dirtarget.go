package fetch

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirTarget string

func (t DirTarget) Open() (TargetInstance, error) {
	return &dirTargetInstance{t, 0}, nil
}

func (t DirTarget) path(name string) string {
	return filepath.Join(string(t), name)
}

type dirTargetInstance struct {
	dest  DirTarget
	count int
}

func (d *dirTargetInstance) Add(name string, data []byte) error {
	d.count++
	return os.WriteFile(d.dest.path(name), data, 0o644)
}

func (d *dirTargetInstance) Finish() (string, error) {
	return fmt.Sprintf("%s: wrote %d files", d.dest, d.count), nil
}
