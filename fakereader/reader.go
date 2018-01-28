package fakereader

import (
	"fmt"
	"io"

	"github.com/albenik/gofaker"
)

type Reader struct {
	t    gofaker.FailTrigger
	flow []io.Reader
	pos  int
}

func New(t gofaker.FailTrigger, flow ...io.Reader) *Reader {
	return &Reader{t: t, flow: flow}
}

func (r *Reader) Read(p []byte) (int, error) {
	var n int
	var err error

	if r.pos < len(r.flow) {
		fr := r.flow[r.pos]
		r.pos++
		if n, err = fr.Read(p); err == nil {
			return n, nil
		}
	} else {
		err = &gofaker.ErrTestFailed{Msg: fmt.Sprintf("unexpected %d read", r.pos+1)}
	}

	if fail, ok := err.(*gofaker.ErrTestFailed); ok {
		err = nil
		r.t.Fatal(fail.Msg)
	}
	return n, err
}
