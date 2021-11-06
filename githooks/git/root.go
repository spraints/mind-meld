package git

import (
	"os/exec"
	"strings"
)

func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	dbg(cmd)
	path, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}
