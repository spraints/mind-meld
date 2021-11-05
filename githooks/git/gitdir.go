package git

import (
	"os/exec"
	"strings"
)

func GetGitDir() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}
