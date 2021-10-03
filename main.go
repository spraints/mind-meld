package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/spraints/mind-meld/lmsp"
	"github.com/spraints/mind-meld/ui"
)

func main() {
	switch len(os.Args) {
	case 1:
		finish(ui.Run())
	case 2:
		finish(dump(os.Args[1]))
	default:
		fmt.Printf("Usage: %s [FILE]\n", os.Args[0])
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
	man, err := l.Manifest()
	if err != nil {
		return err
	}
	spew.Dump(man)

	proj, err := l.Project()
	if err != nil {
		return err
	}
	spew.Dump(proj)

	log.Print("writing JSON back out to 'testing.json'...")
	f, err = os.Create("testing.json")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(proj); err != nil {
		return err
	}

	// todo later - print out programs in pybricks

	return nil
}
