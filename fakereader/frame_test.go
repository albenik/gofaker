package fakereader_test

import (
	"bytes"
	"testing"

	"github.com/albenik/gofaker/fakereader"
	"github.com/stretchr/testify/assert"
)

func TestFrame_Chunk(t *testing.T) {
	data := [][]byte{{1, 2}, {3, 4, 5}, {6, 7, 8, 9}}
	f := fakereader.NewFrame(bytes.Join(data, []byte{}))
	r := fakereader.New("test",
		f.Chunk(2),
		f.Chunk(3),
		f.Chunk(4),
	)

	for _, chunk := range data {
		p := make([]byte, len(chunk))
		n, err := r.Read(p)
		if assert.NoError(t, err) {
			assert.Equal(t, len(chunk), n)
			assert.Equal(t, chunk, p)
		}
	}
}
