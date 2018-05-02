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
		return 0, &gofaker.AssertionFailedError{
			Message: fmt.Sprintf("%s write #%d: unexpected [% X]", w.name, w.wnum+1, p),
		}
	}

	op := w.writers[w.wnum]
	w.wnum++
	n, err := op.Write(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.AssertionFailedError); ok {
			switch top := op.(type) {
			case *WriteOperation:
				err = &gofaker.AssertionFailedError{
					Message: fmt.Sprintf("%s write #%d @ %s:%d: %s", w.name, w.wnum, top.File, top.Line, fail.Message),
				}
			default:
				err = &gofaker.AssertionFailedError{
					Message: fmt.Sprintf("%s write #%d <%#v>: %s", w.name, w.wnum, op, fail.Message),
				}
			}
		}
	}
	return n, err
}

func (w *Writer) Reset() {
	w.wnum = 0
}
