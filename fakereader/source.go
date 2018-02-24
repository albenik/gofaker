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

func AnyData(src []byte) io.Reader {
	return ReaderFunc(func(dst []byte) (int, error) {
		return copy(dst, src), nil
	})
}

func NotLessData(src []byte) io.Reader {
	return ReaderFunc(func(dst []byte) (int, error) {
		if len(dst) < len(src) {
			return 0, &gofaker.ErrTestFailed{Msg: fmt.Sprintf("wrong destination size %d (%d expected)", len(dst), len(src))}
		}
		return copy(dst, src), nil
	})
}

func EqualData(src []byte) io.Reader {
	return ReaderFunc(func(dst []byte) (int, error) {
		if len(dst) != len(src) {
			return 0, &gofaker.ErrTestFailed{Msg: fmt.Sprintf("wrong destination size %d (%d expected)", len(dst), len(src))}
		}
		return copy(dst, src), nil
	})
}

func DelayRead(d time.Duration, r io.Reader, clock *clock.Source) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		clock.Sleep(d)
		return r.Read(p)
	})
}

func ReturnError(err error) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		return 0, err
	})
}
