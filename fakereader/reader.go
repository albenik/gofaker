package fakereader

import (
	"io"

	"github.com/pkg/errors"

	"github.com/albenik/gofaker"
)

type Reader struct {
	ft      gofaker.Fatality
	name    string
	readers []io.Reader
	rnum    int
}

func New(f gofaker.Fatality, n string, flow ...io.Reader) *Reader {
	return &Reader{ft: f, name: n, readers: flow}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.rnum >= len(r.readers) {
		r.ft.Fatalf("%s read #%d: unexpected", r.name, r.rnum+1)
		return 0, nil
	}

	op := r.readers[r.rnum]
	r.rnum++
	n, err := op.Read(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.AssertionFailedError); ok {
			switch top := op.(type) {
			case *ReadOperation:
				r.ft.Fatalf("%s read #%d @ %s:%d: %s", r.name, r.rnum, top.File, top.Line, fail.Message)
			default:
				r.ft.Fatalf("%s read #%d <%#v>: %s", r.name, r.rnum, op, fail.Message)
			}
			return 0, nil
		}
	}
	return n, err
}
