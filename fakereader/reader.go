package fakereader

import (
	"io"

	"github.com/albenik/gofaker"
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

	fr := r.readers[r.rnum]
	r.rnum++
	n, err := fr.Read(p)
	if err != nil {
		if fail, ok := err.(*gofaker.ErrTestFailed); ok {
			err = nil
			r.rnum = len(r.readers) // all next read will fail
			r.t.Fatalf("%s read #%d: %s", r.name, r.rnum, fail.Msg)
		}
	}
	return n, err
}
