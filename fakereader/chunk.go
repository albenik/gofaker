package fakereader

import (
	"fmt"
	"time"
)

type DataChunk struct {
	Bytes  []byte
	Strict bool
	Delay  time.Duration
	Err    error
}

func (c *DataChunk) Read(buf []byte) (int, error) {
	if c.Delay > 0 {
		time.Sleep(c.Delay)
	}
	if len(buf) < len(c.Bytes) {
		return 0, fmt.Errorf("buffer too small: required size %d but provided %d", len(c.Bytes), len(buf))
	}
	if c.Strict && len(buf) > len(c.Bytes) {
		return 0, fmt.Errorf("buffer to large: required size %d but brovides %d", len(c.Bytes), len(buf))
	}
	n := copy(buf, c.Bytes)
	return n, c.Err
}
