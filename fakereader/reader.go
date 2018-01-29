package fakereader

import (
	"io"

	"github.com/albenik/gofaker"
)

type Reader struct {
	t       gofaker.FailTrigger
	readers []io.Reader
	rnum    int
	locked  bool
}

func New(t gofaker.FailTrigger, flow ...io.Reader) *Reader {
	return &Reader{t: t, readers: flow}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.locked {
		r.t.Fatalf("reader locked at %d read", r.rnum+1)
		return 0, nil
	}
	if r.rnum >= len(r.readers) {
		r.t.Fatalf("unexpected %d read", r.rnum+1)
		return 0, nil
	}

	fr := r.readers[r.rnum]
	r.rnum++
	n, err := fr.Read(p)
	if err != nil {
		if fail, ok := err.(*gofaker.ErrTestFailed); ok {
			err = nil
			r.rnum = len(r.readers) // all next read will fail
			r.t.Fatal(fail.Msg)
		}
	}
	return n, err
}
