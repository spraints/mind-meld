package lmsp

import "os"

type File struct {
	f *os.File
}

func Open(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &File{f: f}, nil
}

func (f *File) Close() error {
	return f.f.Close()
}
