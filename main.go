package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spraints/mind-meld/githooks"
	"github.com/spraints/mind-meld/lmsdump"
	"github.com/spraints/mind-meld/lmsp"
	"github.com/spraints/mind-meld/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [dump FILE | browse]\n", os.Args[0])
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
	case "browse":
		finish(ui.Run())
	default:
		fmt.Printf("Usage: %s [dump FILE]\n", os.Args[0])
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
