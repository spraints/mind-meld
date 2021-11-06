package githooks

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spraints/mind-meld/githooks/git"
	"github.com/spraints/mind-meld/lmsdump"
	"github.com/spraints/mind-meld/lmsp/lmspsimple"
)

type PreCommitMode int

const (
	UpdateWorkingCopy PreCommitMode = iota
	UpdateCache
)

func RunPreCommit(mode PreCommitMode) error {
	root, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	files, err := git.LsFiles(root)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Path, ".lms") {
			dumpPath := file.Path + ".dump"
			if err := createDumpFile(root, file, dumpPath, mode); err != nil {
				fmt.Printf("%s: %v\n", dumpPath, err)
			} else {
				fmt.Printf("%s: OK\n", dumpPath)
			}
		}
	}

	return nil
}

func createDumpFile(repoRoot string, file git.IndexEntry, dumpPath string, mode PreCommitMode) error {
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

	switch mode {
	case UpdateWorkingCopy:
		dumped, err := os.Create(filepath.Join(repoRoot, dumpPath))
		if err != nil {
			return err
		}

		if err := lmsdump.Dump(dumped, lmsfile.Project); err != nil {
			return err
		}

		dumped.Close()

		return git.Add(repoRoot, dumpPath)

	case UpdateCache:
		dumpOID, err := git.HashObject(repoRoot, func(w io.Writer) error {
			return lmsdump.Dump(w, lmsfile.Project)
		})
		if err != nil {
			return err
		}

		return git.UpdateIndex(repoRoot, 0100644, dumpOID, dumpPath)

	default:
		return fmt.Errorf("unknown mode %d", mode)
	}
}
