package lmsp

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Reader struct {
	// zr contains three files:
	// - manifest.json
	// - scratch.sb3
	// - icon.svg
	// scratch.sb3 contains several files. one example i have has these:
	// - project.json
	// - 14d134f088239ac481523b3c2c6ecd8c.svg
	// - 93ca32a536da1698ea979f183679af29.png
	zr *zip.Reader
}

var ErrNoManifest = errors.New("no manifest found")

func ReadFile(f *os.File) (*Reader, error) {
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return Read(f, st.Size())
}

func Read(r io.ReaderAt, size int64) (*Reader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	return &Reader{zr: zr}, nil
}

func (r *Reader) Manifest() (Manifest, error) {
	var res Manifest

	f := r.get("manifest.json")
	if f == nil {
		return res, ErrNoManifest
	}

	fr, err := f.Open()
	if err != nil {
		return res, err
	}
	defer fr.Close()

	err = json.NewDecoder(fr).Decode(&res)
	return res, err
}

func (r *Reader) get(name string) *zip.File {
	for _, f := range r.zr.File {
		if f.Name == name {
			return f
		}
	}
	return nil
}
