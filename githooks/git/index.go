package git

import "os/exec"

func Add(cwd, path string) error {
	cmd := exec.Command("git", "add", path)
	cmd.Dir = cwd
	_, err := cmd.Output()
	return err
}
