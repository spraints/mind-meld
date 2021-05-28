package lmsp

import (
	"os"
	"testing"
)

func TestReadingFile(t *testing.T) {
	run := func(filename, desc string, tf func(*testing.T, *File)) {
		t.Run(filename+": "+desc, func(t *testing.T) {
			t.Parallel()

			f, err := Open(filename)
			if err != nil {
				t.Errorf("error opening file: %v", err)
				return
			}
			defer f.Close()
			tf(t, f)
		})
	}

	t.Run("no file", func(t *testing.T) {
		_, err := Open("does-not-exist.lmsp")
		if !os.IsNotExist(err) {
			t.Errorf("expected IsNotExist error but got %v", err)
		}
	})

	run("testdata/my-block-with-no-refs.lmsp", "stub", func(t *testing.T, f *File) {
		// todo
	})

	run("testdata/Gyro drive.lmsp", "stub", func(t *testing.T, f *File) {
		// todo
	})
}
