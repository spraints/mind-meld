package lmspsimple

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/spraints/mind-meld/lmsp"
)

type File struct {
	Raw     []byte
	Project lmsp.Project
}

func Read(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(f, st.Size())
	if err != nil {
		return nil, err
	}

	zf, err := zr.Open("scratch.sb3")
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	return readScratch(zf)
}

func readScratch(f fs.File) (*File, error) {
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}

	project, err := zr.Open("project.json")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(project)
	if err != nil {
		return nil, err
	}

	res := &File{Raw: data}
	if err := json.Unmarshal(data, &res.Project); err != nil {
		return nil, err
	}

	return res, nil
}
