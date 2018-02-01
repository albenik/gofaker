package fakewriter

import (
	"io"

	"github.com/albenik/gofaker"
)

type Writer struct {
	t       gofaker.FailTrigger
	writers []io.Writer
	wnum    int
	locked  bool
}

func New(t gofaker.FailTrigger, flow ...io.Writer) *Writer {
	return &Writer{t: t, writers: flow}
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.locked {
		w.t.Fatalf("writer locked at %d write [% X]", w.wnum+1, p)
		return 0, nil
	}
	if w.wnum >= len(w.writers) {
		w.t.Fatalf("unexpected %d write [% X]", w.wnum+1, p)
		return 0, nil
	}

	fw := w.writers[w.wnum]
	w.wnum++
	n, err := fw.Write(p)
	if err != nil {
		if fail, ok := err.(*gofaker.ErrTestFailed); ok {
			w.locked = true
			w.t.Fatal(fail.Msg)
			return n, nil
		}
	}
	return n, err
}

func (w *Writer) Reset() {
	w.wnum = 0
}
