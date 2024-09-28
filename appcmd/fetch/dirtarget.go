package fetch

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirTarget string

func (t DirTarget) Open() (TargetInstance, error) {
	if st, err := os.Stat(string(t)); err != nil {
		return nil, err
	} else if !st.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", t)
	}
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
	destFile := d.dest.path(name)
	if err := os.MkdirAll(filepath.Dir(destFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destFile, data, 0o644)
}

func (d *dirTargetInstance) Finish() (string, error) {
	return fmt.Sprintf("%s: wrote %d files", d.dest, d.count), nil
}
