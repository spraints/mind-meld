package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spraints/mind-meld/githooks"
	"github.com/spraints/mind-meld/lmsdump"
	"github.com/spraints/mind-meld/lmsp"
	"github.com/spraints/mind-meld/ui"
)

func main() {
	usage := fmt.Sprintf("Usage: %s [dump FILE | ls | browse]", os.Args[0])

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "pre-commit":
		mode := githooks.UpdateWorkingCopy
		for _, arg := range os.Args[2:] {
			if arg == "--cached" {
				mode = githooks.UpdateCache
			}
		}
		finish(githooks.RunPreCommit(mode))
	case "dump":
		finish(dump(os.Args[2]))
	case "ls":
		finish(ls())
	case "browse":
		finish(ui.Run())
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}

func finish(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func dump(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	l, err := lmsp.ReadFile(f)
	if err != nil {
		return err
	}
	/*
		man, err := l.Manifest()
		if err != nil {
			return err
		}
		spew.Dump(man)
	*/

	proj, err := l.Project()
	if err != nil {
		return err
	}
	lmsdump.Dump(os.Stdout, proj)

	if os.Getenv("WRITE_PROJECT_JSON") != "" {
		log.Print("writing JSON back out to 'testing.json'...")
		f, err = os.Create("testing.json")
		if err != nil {
			return err
		}
		defer f.Close()
		if err := json.NewEncoder(f).Encode(proj); err != nil {
			return err
		}
	}

	// todo later - print out programs in pybricks

	return nil
}

func ls() error {
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME must be set in order to find mindstorms files")
	}

	const (
		spikenextDir     = "Library/Containers/com.lego.education.spikenext/Data/Documents/LEGO Education SPIKE"
		robotinventorDir = "Library/Containers/com.lego.retail.mindstorms.robotinventor/Data/Documents/LEGO MINDSTORMS"
		ev3Dir           = "Documents/LEGO Education EV3 Content"
	)

	c := exec.Command("find",
		".",
		filepath.Join(home, spikenextDir),
		filepath.Join(home, robotinventorDir),
		filepath.Join(home, ev3Dir),
		"-type", "f",
		"(", "-name", "*.lmsp", "-or", "-name", "*.lms", "-or", "-name", "*.llsp3", ")",
		"-ls")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
