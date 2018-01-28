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

func And(w ...io.Writer) io.Writer {
	return WriterFunc(func(data []byte) (int, error) {
		var (
			n   int
			err error
		)
		for _, ww := range w {
			if n, err = ww.Write(data); err != nil {
				break
			}
		}
		return n, err
	})
}

func ExpectLen(l int) io.Writer {
	return WriterFunc(func(data []byte) (int, error) {
		if len(data) != l {
			return len(data), &gofaker.ErrTestFailed{Msg: fmt.Sprintf("invalid data length: %d expected but %d recieved", l, len(data))}
		}
		return len(data), nil
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

func ShortWrite(l int) io.Writer {
	return WriterFunc(func(data []byte) (int, error) {
		if l < len(data) {
			return l, nil
		}
		return len(data), nil
	})
}

func DelayWrite(d time.Duration, clock *clock.Source) io.Writer {
	return WriterFunc(func(data []byte) (int, error) {
		clock.Sleep(d)
		return len(data), nil
	})
}
