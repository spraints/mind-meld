package lmspsimple

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
)

type File struct {
	Names []string
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

	names := make([]string, 0, len(zr.File))
	for _, f := range zr.File {
		names = append(names, f.Name)
	}

	return &File{Names: names}, nil
}
