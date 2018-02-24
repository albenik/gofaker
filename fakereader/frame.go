package fakereader

import "io"

type Frame struct {
	data []byte
	offs int
}

func NewFrame(p []byte) *Frame {
	f := &Frame{data: make([]byte, len(p))}
	copy(f.data, p)
	return f
}

func (f *Frame) Chunk(ln int) io.Reader {
	r := EqualData(f.data[f.offs : f.offs+ln])
	f.offs += ln
	return r
}
