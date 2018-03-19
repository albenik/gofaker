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
		panic(fmt.Sprintf("%s read #%d: unexpected", r.name, r.rnum+1))
	}

	op := r.readers[r.rnum]
	r.rnum++
	n, err := op.Read(p)
	if err != nil {
		if fail, ok := errors.Cause(err).(*gofaker.CheckFailed); ok {
			var msg string

			switch top := op.(type) {
			case *ReadOperation:
				msg = fmt.Sprintf("%s read #%d: %s @ %s:%d", r.name, r.rnum, fail.Message, top.File, top.Line)
			default:
				msg = fmt.Sprintf("%s read #%d: %s <%#v>", r.name, r.rnum, fail.Message, op)
			}
			panic(msg)
		}
	}
	return n, err
}
