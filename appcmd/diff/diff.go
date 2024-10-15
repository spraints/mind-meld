package diff

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spraints/mind-meld/appcmd"
	"github.com/spraints/mind-meld/appcmd/fetch"
)

func Run(app appcmd.App, baseRev string) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	t := &target{
		repo: repo,
		tb:   fetch.NewTreeBuilder(repo),
	}

	if _, err := fetch.Run(app, t); err != nil {
		return err
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	args := []string{
		"git",
		"diff",
		baseRev,
		t.newTreeID.String(),
	}
	if os.Getenv("DEBUG") == "1" {
		fmt.Printf("%s\n", strings.Join(args, " "))
	}
	return syscall.Exec(gitPath, args, os.Environ())
}

type target struct {
	repo *git.Repository
	tb   *fetch.TreeBuilder

	newTreeID *plumbing.Hash
}

func (t *target) Open() (fetch.TargetInstance, error) {
	return t, nil
}

func (t *target) PathSeparator() string {
	return fetch.GitPathSeparator
}

func (t *target) Add(name string, data []byte) error {
	return t.tb.Add(name, data)
}

func (t *target) Finish() (string, error) {
	treeID, err := t.tb.Finish()
	if err != nil {
		return "", err
	}

	t.newTreeID = &treeID
	return "", nil
}
