package fakewriter

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/albenik/gofaker"
	"github.com/albenik/gofaker/clock"
)

type WriterFunc func([]byte) (int, error)

func (f WriterFunc) Write(p []byte) (int, error) {
	return f(p)
}

func ExpectLen(exp int) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		if len(p) != exp {
			return len(p), &gofaker.ErrTestFailed{Msg: fmt.Sprintf("invalid data length %d [% X] (expected %d)", len(p), p, exp)}
		}
		return len(p), nil
	})
}

func ExpectData(exp []byte) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		if !bytes.Equal(exp, p) {
			return len(p), &gofaker.ErrTestFailed{Msg: fmt.Sprintf("invalid data [% X] (expected [% X])", p, exp)}
		}
		return len(p), nil
	})
}

func ShortWrite(l int, w io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		n, err := w.Write(p)
		if err != nil {
			return n, err
		}
		if l < len(p) {
			return l, nil
		}
		return len(p), nil
	})
}

func DelayWrite(d time.Duration, w io.Writer, clock *clock.Source) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		clock.Sleep(d)
		return w.Write(p)
	})
}

func ForceLen(l int, w io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		_, err := w.Write(p)
		return l, err
	})
}

func WriteError(err error) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		return 0, err
	})
}
