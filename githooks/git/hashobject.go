package git

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func HashObject(repoRoot string, writeStdin func(io.Writer) error) (string, error) {
	cmd := exec.Command("git", "hash-object", "-w", "--stdin")
	dbg(cmd)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	waitErr := make(chan error)
	go func() {
		defer close(waitErr)
		waitErr <- cmd.Wait()
	}()

	writeStdin(stdin)
	stdin.Close()
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), <-waitErr
}
