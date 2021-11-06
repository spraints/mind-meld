package git

import (
	"os/exec"
	"strings"
)

func GetGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	dbg(cmd)
	path, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}
