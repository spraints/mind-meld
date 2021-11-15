package lmsdump

import (
	"bytes"
	"io"
)

func indent(w io.Writer) io.Writer {
	return &indented{w: w, i: false}
}

func indentStartingNow(w io.Writer) io.Writer {
	return &indented{w: w, i: true}
}

type indented struct {
	w io.Writer
	i bool
}

var indentPadding = []byte("  ")

func (i *indented) Write(p []byte) (int, error) {
	n := 0
	for len(p) > 0 {
		if i.i {
			_, err := i.w.Write(indentPadding)
			if err != nil {
				return n, err
			}
			i.i = false
		}
		x := bytes.IndexRune(p, '\n')
		if x == -1 {
			nn, err := i.w.Write(p)
			n += nn
			return n, err
		}
		nn, err := i.w.Write(p[:x+1])
		n += nn
		if err != nil {
			return n, err
		}
		i.i = true
		p = p[x+1:]
	}
	return n, nil
}
