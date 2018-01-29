package fakereader

import (
	"fmt"
	"io"
	"time"

	"github.com/albenik/gofaker"
	"github.com/albenik/gofaker/clock"
)

type ReaderFunc func([]byte) (int, error)

func (f ReaderFunc) Read(p []byte) (int, error) {
	return f(p)
}

func StrictBytesReader(data []byte) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		n := copy(p, data)
		if n != len(data) {
			return n, &gofaker.ErrTestFailed{Msg: fmt.Sprintf("expected buffer length is %d but actual is %d", len(data), n)}
		}
		return n, nil
	})
}

func DelayRead(d time.Duration, r io.Reader, clock *clock.Source) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		clock.Sleep(d)
		return r.Read(p)
	})
}
