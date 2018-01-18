package fakewriter

import (
	"bytes"
	"fmt"
	"io"
)

type WriterFunc func([]byte) (int, error)

func (f WriterFunc) Write(data []byte) (int, error) {
	return f(data)
}

func AssertLen(l int) io.Writer {
	return WriterFunc(func(data []byte) (int, error) {
		if len(data) != l {
			return 0, fmt.Errorf("invalid data length: %d expected but %d recieved", l, len(data))
		}
		return len(data), nil
	})
}

func AssertData(a []byte) io.Writer {
	return WriterFunc(func(b []byte) (int, error) {
		if !bytes.Equal(a, b) {
			return 0, fmt.Errorf("invalid data: [% X] expected but [% X] recieved", a, b)
		}
		return len(b), nil
	})
}
