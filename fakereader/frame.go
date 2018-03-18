package fakereader

import (
	"fmt"
	"io"

	"github.com/albenik/gofaker"
)

type Frame struct {
	data []byte
	offs int
}

func NewFrame(p []byte) *Frame {
	f := &Frame{data: make([]byte, len(p))}
	copy(f.data, p)
	return f
}

func (f *Frame) Restart() *Frame {
	f.offs = 0
	return f
}

func (f *Frame) Chunk(ln int) io.Reader {
	end := f.offs + ln
	if end > len(f.data) {
		panic("not enougth data")
	}
	src := f.data[f.offs:end]
	file, line := gofaker.GetSourceCodeLine(4)
	r := &ReadOperation{
		File: file,
		Line: line,
		Fn: func(dst []byte) (int, error) {
			if len(dst) != len(src) {
				return 0, &gofaker.ErrTestFailed{Msg: fmt.Sprintf("wrong destination size %d (%d expected)", len(dst), len(src))}
			}
			return copy(dst, src), nil
		},
	}
	f.offs = end
	return r
}

func (f *Frame) MultipleChunks(ln, cnt int) []io.Reader {
	if f.offs+ln*cnt > len(f.data) {
		panic("not enougth data")
	}
	list := make([]io.Reader, cnt)
	for i := 0; i < cnt; i++ {
		list[i] = f.Chunk(ln)
	}
	return list
}

func (f *Frame) AllChunks(ln int) []io.Reader {
	if f.offs >= len(f.data) {
		panic("not enougth data")
	}
	remain := len(f.data) - f.offs
	if remain%ln != 0 {
		panic("data chunks not aligned")
	}
	return f.MultipleChunks(ln, remain/ln)
}
