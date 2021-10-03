package lmspsimple

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
)

type File struct {
	JSON string
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

	json, err := zr.Open("project.json")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(json)
	if err != nil {
		return nil, err
	}

	return &File{JSON: string(data)}, nil
}
