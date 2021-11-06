package git

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type IndexEntry struct {
	Mode  int
	OID   string
	Index int
	Path  string
}

func LsFiles() ([]IndexEntry, error) {
	cmd := exec.Command("git", "ls-files", "-s")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	waitErr := make(chan error)
	go func() {
		defer close(waitErr)
		waitErr <- cmd.Wait()
	}()
	entries, err := parseLsFiles(stdout)
	if err != nil {
		return nil, err
	}
	return entries, <-waitErr
}

func parseLsFiles(r io.Reader) ([]IndexEntry, error) {
	br := bufio.NewReader(r)
	var entries []IndexEntry
	for {
		entry, err := parseLsFilesLine(br)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func parseLsFilesLine(r *bufio.Reader) (IndexEntry, error) {
	modeStr, err := r.ReadString(' ')
	if err != nil {
		return IndexEntry{}, err
	}
	mode, err := strconv.ParseInt(strings.TrimSpace(modeStr), 8, 32)
	if err != nil {
		return IndexEntry{}, err
	}
	oid, err := r.ReadString(' ')
	if err != nil {
		return IndexEntry{}, err
	}
	indexStr, err := r.ReadString('\t')
	if err != nil {
		return IndexEntry{}, err
	}
	index, err := strconv.ParseInt(strings.TrimSpace(indexStr), 10, 32)
	if err != nil {
		return IndexEntry{}, err
	}
	path, err := r.ReadString('\n')
	if err != nil {
		return IndexEntry{}, err
	}
	return IndexEntry{
		Mode:  int(mode),
		OID:   strings.TrimSpace(oid),
		Index: int(index),
		Path:  strings.TrimSpace(path),
	}, nil
}
