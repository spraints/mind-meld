package watch

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spraints/mind-meld/appcmd"
	"github.com/spraints/mind-meld/appcmd/fetch"
)

func Run(ctx context.Context, a appcmd.App, t fetch.Target) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	anyOK := false
	for _, d := range a.ProjectDirs() {
		if err := watcher.Add(d); err != nil {
			fmt.Printf("%s: %v\n", d, err)
		} else {
			fmt.Printf("watching %s\n", d)
			anyOK = true
		}
	}

	if !anyOK {
		return fmt.Errorf("no watchable directories found")
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case _, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if err := fetch.Run(a, t); err != nil {
				return err
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return err
		}
	}
}
