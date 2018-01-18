package fakereader

import (
	"fmt"
	"io"
)

type Reader struct {
	Flow []io.Reader
	pos  int
}

func (sr *Reader) Read(buf []byte) (int, error) {
	if sr.pos < len(sr.Flow) {
		fr := sr.Flow[sr.pos]
		sr.pos++
		return fr.Read(buf)
	}
	return 0, fmt.Errorf("unexpected %d read", sr.pos+1)
}
