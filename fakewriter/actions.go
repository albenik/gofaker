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

func ExpectLen(l int) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		if len(p) != l {
			return len(p), &gofaker.ErrTestFailed{Msg: fmt.Sprintf("invalid data length: %d expected but %d recieved", l, len(p))}
		}
		return len(p), nil
	})
}

func ExpectData(a []byte) io.Writer {
	return WriterFunc(func(b []byte) (int, error) {
		if !bytes.Equal(a, b) {
			return len(b), &gofaker.ErrTestFailed{Msg: fmt.Sprintf("invalid data: [% X] expected but [% X] recieved", a, b)}
		}
		return len(b), nil
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
