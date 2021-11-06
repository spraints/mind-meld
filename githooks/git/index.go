package git

import (
	"fmt"
	"os/exec"
)

func Add(cwd, path string) error {
	cmd := exec.Command("git", "add", path)
	dbg(cmd)
	cmd.Dir = cwd
	_, err := cmd.Output()
	return err
}

func UpdateIndex(cwd string, mode int, oid, path string) error {
	cmd := exec.Command("git", "update-index", "--add", "--cacheinfo", indexItem(mode, oid, path))
	dbg(cmd)
	cmd.Dir = cwd
	_, err := cmd.Output()
	return err
}

func indexItem(mode int, oid, path string) string {
	return fmt.Sprintf("%o,%s,%s", mode, oid, path)
}
