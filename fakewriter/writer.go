package fakewriter

import (
	"io"

	"github.com/pkg/errors"

	"github.com/albenik/gofaker"
)

type Writer struct {
	ft      gofaker.Fatality
	name    string
	writers []io.Writer
	wnum    int
}

func New(f gofaker.Fatality, n string, flow ...io.Writer) *Writer {
	return &Writer{ft: f, name: n, writers: flow}
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.wnum >= len(w.writers) {
		w.ft.Fatalf("%s write #%d: unexpected [% X]", w.name, w.wnum+1, p)
		return 0, nil
	}

	op := w.writers[w.wnum]
	w.wnum++
	n, err := op.Write(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.AssertionFailedError); ok {
			switch top := op.(type) {
			case *WriteOperation:
				w.ft.Fatalf("%s write #%d @ %s:%d: %s", w.name, w.wnum, top.File, top.Line, fail.Message)
			default:
				w.ft.Fatalf("%s write #%d <%#v>: %s", w.name, w.wnum, op, fail.Message)
			}
			return 0, nil
		}
	}
	return n, err
}

func (w *Writer) Reset() {
	w.wnum = 0
}
