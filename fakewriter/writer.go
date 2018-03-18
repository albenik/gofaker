package fakewriter

import (
	"io"

	"github.com/albenik/gofaker"
	"github.com/pkg/errors"
)

type Writer struct {
	t       gofaker.FailTrigger
	name    string
	writers []io.Writer
	wnum    int
	locked  bool
}

func New(t gofaker.FailTrigger, n string, flow ...io.Writer) *Writer {
	return &Writer{t: t, name: n, writers: flow}
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.locked {
		w.t.Fatalf("%s write #%d: locked [% X]", w.name, w.wnum+1, p)
		return 0, nil
	}
	if w.wnum >= len(w.writers) {
		w.t.Fatalf("%s write #%d: unexpected [% X]", w.name, w.wnum+1, p)
		return 0, nil
	}

	op := w.writers[w.wnum]
	w.wnum++
	n, err := op.Write(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.ErrTestFailed); ok {
			w.locked = true
			switch top := op.(type) {
			case *WriteOperation:
				w.t.Fatalf("%s write #%d: %s @ %s:%d", w.name, w.wnum, fail.Msg, top.File, top.Line)
			default:
				w.t.Fatalf("%s write #%d: %s <%#v>", w.name, w.wnum, fail.Msg, op)
			}
			return n, nil
		}
	}
	return n, err
}

func (w *Writer) Reset() {
	w.wnum = 0
}
