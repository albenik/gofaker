package fakewriter

import (
	"fmt"
	"io"
)

type Writer struct {
	Flow []io.Writer
	pos  int
}

func (w *Writer) Write(data []byte) (int, error) {
	if w.pos < len(w.Flow) {
		fw := w.Flow[w.pos]
		w.pos++
		return fw.Write(data)
	}
	return 0, fmt.Errorf("unexpected %d write", w.pos+1)
}
