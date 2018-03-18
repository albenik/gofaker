package fakereader

import (
	"fmt"
	"io"
	"time"

	"github.com/albenik/gofaker"
	"github.com/albenik/gofaker/clock"
)

type ReadOperation struct {
	File string
	Line int
	Fn   func([]byte) (int, error)
}

func (ro *ReadOperation) Read(p []byte) (int, error) {
	return ro.Fn(p)
}

func AnyData(src []byte) io.Reader {
	file, line := gofaker.GetSourceCodeLine(3)
	return &ReadOperation{
		File: file,
		Line: line,
		Fn: func(dst []byte) (int, error) {
			return copy(dst, src), nil
		},
	}
}

func NotLessData(src []byte) io.Reader {
	file, line := gofaker.GetSourceCodeLine(3)
	return &ReadOperation{
		File: file,
		Line: line,
		Fn: func(dst []byte) (int, error) {
			if len(dst) < len(src) {
				return 0, &gofaker.ErrTestFailed{Msg: fmt.Sprintf("wrong destination size %d (%d expected)", len(dst), len(src))}
			}
			return copy(dst, src), nil
		},
	}
}

func EqualData(src []byte) io.Reader {
	file, line := gofaker.GetSourceCodeLine(3)
	return &ReadOperation{
		File: file,
		Line: line,
		Fn: func(dst []byte) (int, error) {
			if len(dst) != len(src) {
				return 0, &gofaker.ErrTestFailed{Msg: fmt.Sprintf("wrong destination size %d (%d expected)", len(dst), len(src))}
			}
			return copy(dst, src), nil
		},
	}
}

func DelayRead(d time.Duration, r io.Reader, clock *clock.Source) io.Reader {
	file, line := gofaker.GetSourceCodeLine(3)
	return &ReadOperation{
		File: file,
		Line: line,
		Fn: func(p []byte) (int, error) {
			clock.Sleep(d)
			return r.Read(p)
		},
	}
}

func ReturnError(err error) io.Reader {
	file, line := gofaker.GetSourceCodeLine(3)
	return &ReadOperation{
		File: file,
		Line: line,
		Fn: func(p []byte) (int, error) {
			return 0, err
		},
	}
}
