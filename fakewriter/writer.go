package fakewriter

import (
	"fmt"
	"io"

	"github.com/albenik/gofaker"
)

type Writer struct {
	t    gofaker.FailTrigger
	flow []io.Writer
	pos  int
}

func New(t gofaker.FailTrigger, flow ...io.Writer) *Writer {
	return &Writer{t: t, flow: flow}
}

func (w *Writer) Write(p []byte) (int, error) {
	var n int
	var err error

	if w.pos < len(w.flow) {
		fw := w.flow[w.pos]
		w.pos++
		if n, err = fw.Write(p); err == nil {
			return n, nil
		}
	} else {
		err = &gofaker.ErrTestFailed{Msg: fmt.Sprintf("unexpected %d write", w.pos+1)}
	}

	if fail, ok := err.(*gofaker.ErrTestFailed); ok {
		err = nil
		w.t.Fatal(fail.Msg)
	}
	return n, err
}
