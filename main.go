package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spraints/mind-meld/lmsp"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatal(err)
	}
	l, err := lmsp.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	man, err := l.Manifest()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(man)

	proj, err := l.Project()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(proj)

	log.Print("writing JSON back out to 'testing.json'...")
	f, err = os.Create("testing.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(proj); err != nil {
		log.Fatal(err)
	}

	// todo later - print out programs in pybricks
}
