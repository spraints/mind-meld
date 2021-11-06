package githooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spraints/mind-meld/githooks/git"
	"github.com/spraints/mind-meld/lmsdump"
	"github.com/spraints/mind-meld/lmsp/lmspsimple"
)

func RunPreCommit() error {
	root, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	files, err := git.LsFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Path, ".lms") {
			dumpPath := file.Path + ".dump"
			if err := createDumpFile(root, file, dumpPath); err != nil {
				fmt.Printf("%s: %v\n", dumpPath, err)
			} else {
				fmt.Printf("%s: OK\n", dumpPath)
			}
		}
	}

	return nil
}

func createDumpFile(repoRoot string, file git.IndexEntry, dumpPath string) error {
	tmp, err := os.CreateTemp("", "mind-meld-pre-commit-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := git.ReadObjectToFile(file.OID, tmp); err != nil {
		return fmt.Errorf("read blob %s: %w", file.OID[:10], err)
	}

	tmp.Close()

	lmsfile, err := lmspsimple.Read(tmp.Name())
	if err != nil {
		return err
	}

	dumped, err := os.Create(filepath.Join(repoRoot, dumpPath))
	if err != nil {
		return err
	}

	if err := lmsdump.Dump(dumped, lmsfile.Project); err != nil {
		return err
	}

	dumped.Close()

	return git.Add(repoRoot, dumpPath)
}
