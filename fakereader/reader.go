package fakereader

import (
	"fmt"
	"io"

	"github.com/albenik/gofaker"
	"github.com/pkg/errors"
)

type Reader struct {
	name    string
	readers []io.Reader
	rnum    int
}

func New(n string, flow ...io.Reader) *Reader {
	return &Reader{name: n, readers: flow}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.rnum >= len(r.readers) {
		return 0, &gofaker.AssertionFailedError{
			Message: fmt.Sprintf("%s read #%d: unexpected", r.name, r.rnum+1),
		}
	}

	op := r.readers[r.rnum]
	r.rnum++
	n, err := op.Read(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.AssertionFailedError); ok {
			switch top := op.(type) {
			case *ReadOperation:
				err = &gofaker.AssertionFailedError{
					Message: fmt.Sprintf("%s read #%d @ %s:%d: %s", r.name, r.rnum, top.File, top.Line, fail.Message),
				}
			default:
				err = &gofaker.AssertionFailedError{
					Message: fmt.Sprintf("%s read #%d <%#v>: %s", r.name, r.rnum, op, fail.Message),
				}
			}
		}
	}
	return n, err
}
