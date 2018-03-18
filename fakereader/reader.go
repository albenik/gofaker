package fakereader

import (
	"io"

	"github.com/albenik/gofaker"
	"github.com/pkg/errors"
)

type Reader struct {
	t       gofaker.FailTrigger
	name    string
	readers []io.Reader
	rnum    int
	locked  bool
}

func New(t gofaker.FailTrigger, n string, flow ...io.Reader) *Reader {
	return &Reader{t: t, name: n, readers: flow}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.locked {
		r.t.Fatalf("%s read #%d: reader locked", r.name, r.rnum+1)
		return 0, nil
	}
	if r.rnum >= len(r.readers) {
		r.t.Fatalf("%s read #%d: unexpected", r.name, r.rnum+1)
		return 0, nil
	}

	op := r.readers[r.rnum]
	r.rnum++
	n, err := op.Read(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.ErrTestFailed); ok {
			r.locked = true
			switch top := op.(type) {
			case *ReadOperation:
				r.t.Fatalf("%s read #%d: %s @ %s:%d", r.name, r.rnum, fail.Msg, top.File, top.Line)
			default:
				r.t.Fatalf("%s read #%d: %s <%#v>", r.name, r.rnum, fail.Msg, op)
			}
			return n, nil
		}
	}
	return n, err
}
