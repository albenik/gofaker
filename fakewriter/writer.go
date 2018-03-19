package fakewriter

import (
	"fmt"
	"io"

	"github.com/albenik/gofaker"
	"github.com/pkg/errors"
)

type Writer struct {
	name    string
	writers []io.Writer
	wnum    int
}

func New(n string, flow ...io.Writer) *Writer {
	return &Writer{name: n, writers: flow}
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.wnum >= len(w.writers) {
		panic(fmt.Sprintf("%s write #%d: unexpected [% X]", w.name, w.wnum+1, p))
	}

	op := w.writers[w.wnum]
	w.wnum++
	n, err := op.Write(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.CheckFailed); ok {
			var msg string
			switch top := op.(type) {
			case *WriteOperation:
				msg = fmt.Sprintf("%s write #%d: %s @ %s:%d", w.name, w.wnum, fail.Message, top.File, top.Line)
			default:
				msg = fmt.Sprintf("%s write #%d: %s <%#v>", w.name, w.wnum, fail.Message, op)
			}
			panic(msg)
		}
	}
	return n, err
}

func (w *Writer) Reset() {
	w.wnum = 0
}
