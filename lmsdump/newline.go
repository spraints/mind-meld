package lmsdump

import (
	"fmt"
	"io"
)

func finishWithNewline(w io.Writer) *finisher {
	return &finisher{w: w}
}

type finisher struct {
	w       io.Writer
	needsNL bool
}

func (f *finisher) Write(data []byte) (int, error) {
	if len(data) > 0 {
		f.needsNL = data[len(data)-1] != '\n'
	}
	return f.w.Write(data)
}

func (f *finisher) Finish() {
	if f.needsNL {
		fmt.Fprintln(f.w)
		f.needsNL = false
	}
}
