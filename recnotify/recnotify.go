package recnotify

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func MaybeAddRecursive(w *fsnotify.Watcher, evt fsnotify.Event) error {
	if st, err := os.Stat(evt.Name); err == nil && st.IsDir() {
		if err := AddRecursive(w, evt.Name); err != nil {
			return err
		}
	}

	return nil
}

func AddRecursive(w *fsnotify.Watcher, path string) error {
	if err := filepath.WalkDir(path, recWatchFn(w)); err != nil {
		return err
	}

	return nil
}

func recWatchFn(w *fsnotify.Watcher) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// There was a problem reading path. Just skip this one entry silently for now.
			return nil
		}
		if d.IsDir() {
			if err := w.Add(path); err != nil {
				return err
			}
		}
		return nil
	}
}
