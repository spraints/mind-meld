package main

import (
	"log"
	"os"

	"github.com/spraints/mind-meld/lmsp"
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
        log.Printf("%#v\n", man)
}
