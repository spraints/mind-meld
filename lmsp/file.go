package lmsp

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
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

var (
	ErrNoManifest = errors.New("no manifest found")
	ErrNoScratch  = errors.New("no scratch data found")
	ErrNoProject  = errors.New("no project found")
	ErrNoPython   = errors.New("no python code found")
)

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

// Manifest reads the manifest header information from the file.
func (r *Reader) Manifest() (Manifest, error) {
	var res Manifest

	f := get(r.zr, "manifest.json")
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

// Python reads the python source from projectbody.json in the lego zip file.
func (r *Reader) Python() (map[string]string, error) {
	f := get(r.zr, "projectbody.json")
	if f == nil {
		return nil, ErrNoPython
	}

	fr, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	res := map[string]string{}
	err = json.NewDecoder(fr).Decode(&res)
	return res, err
}

// Project reads the scratch project and programs from the file.
func (r *Reader) Project() (Project, error) {
	var res Project

	f := get(r.zr, "scratch.sb3")
	if f == nil {
		return res, ErrNoScratch
	}

	fr, err := f.Open()
	if err != nil {
		return res, err
	}
	defer fr.Close()

	data, err := ioutil.ReadAll(fr)
	if err != nil {
		return res, err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return res, err
	}

	f = get(zr, "project.json")
	if f == nil {
		return res, ErrNoProject
	}

	pr, err := f.Open()
	if err != nil {
		return res, err
	}
	defer pr.Close()

	err = json.NewDecoder(pr).Decode(&res)
	return res, err
}

func get(r *zip.Reader, name string) *zip.File {
	for _, f := range r.File {
		if f.Name == name {
			return f
		}
	}
	return nil
}
