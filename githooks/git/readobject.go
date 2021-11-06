package git

import (
	"io"
	"os"
	"os/exec"
)

func ReadObjectToFile(oid string, w io.Writer) error {
	cmd := exec.Command("git", "cat-file", "-p", oid)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
