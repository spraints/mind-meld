package githooks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spraints/mind-meld/githooks/git"
)

const (
	PreCommit = "pre-commit"
)

func Install() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	gitdir, err := git.GetGitDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(gitdir, "hooks", PreCommit)
	if existing, err := os.Readlink(hookPath); err == nil && filepath.Base(existing) == "mind-meld" {
		fmt.Printf("removing previous mind-meld hook link to %s\n", existing)
		if err := os.Remove(hookPath); err != nil {
			return err
		}
	}
	if err := os.Symlink(exe, hookPath); err != nil {
		return err
	}

	fmt.Printf("created %s as link to %s\n", hookPath, exe)
	return nil
}
