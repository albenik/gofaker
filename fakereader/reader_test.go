package fakereader_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/albenik/gofaker/fakereader"
	"github.com/stretchr/testify/assert"
)

func TestReader_Read(t *testing.T) {
	r := fakereader.Reader{Flow: []io.Reader{
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}},
		&fakereader.DataChunk{Bytes: []byte{1, 2}},
		&fakereader.DataChunk{Bytes: []byte{1}},
	}}

	read := func(expect []byte) {
		buf := make([]byte, len(expect))
		n, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, len(expect), n)
		assert.Equal(t, expect, buf)
	}

	read([]byte{1, 2, 3})
	read([]byte{1, 2})
	read([]byte{1})

	n, err := r.Read(make([]byte, 3))
	assert.EqualError(t, err, "unexpected 4 read")
	assert.Equal(t, 0, n)
}

func TestDataChunk_ReadStrict(t *testing.T) {
	r := fakereader.Reader{Flow: []io.Reader{
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true},
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true},
	}}

	n, err := r.Read(make([]byte, 3))
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = r.Read(make([]byte, 7))
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = r.Read(make([]byte, 1))
	assert.EqualError(t, err, "buffer too small: required size 3 but provided 1")
	assert.Equal(t, 0, n)

	n, err = r.Read([]byte{})
	assert.EqualError(t, err, "buffer too small: required size 3 but provided 0")
	assert.Equal(t, 0, n)

	n, err = r.Read(nil)
	assert.EqualError(t, err, "buffer too small: required size 3 but provided 0")
	assert.Equal(t, 0, n)

	n, err = r.Read(make([]byte, 3))
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = r.Read(make([]byte, 7))
	assert.EqualError(t, err, "buffer to large: required size 3 but brovides 7")
	assert.Equal(t, 0, n)

	n, err = r.Read(make([]byte, 1))
	assert.EqualError(t, err, "buffer too small: required size 3 but provided 1")
	assert.Equal(t, 0, n)

	n, err = r.Read([]byte{})
	assert.EqualError(t, err, "buffer too small: required size 3 but provided 0")
	assert.Equal(t, 0, n)

	n, err = r.Read(nil)
	assert.EqualError(t, err, "buffer too small: required size 3 but provided 0")
	assert.Equal(t, 0, n)
}

func TestDataChunk_ReadWithErr(t *testing.T) {
	r := fakereader.Reader{Flow: []io.Reader{
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Err: errors.New("test error 1")},
		&fakereader.DataChunk{Bytes: []byte{}, Err: errors.New("test error 2")},
		&fakereader.DataChunk{Bytes: nil, Err: errors.New("test error 3")},
	}}

	n, err := r.Read(make([]byte, 3))
	assert.EqualError(t, err, "test error 1")
	assert.Equal(t, 3, n)

	n, err = r.Read(make([]byte, 3))
	assert.EqualError(t, err, "test error 2")
	assert.Equal(t, 0, n)

	n, err = r.Read(make([]byte, 3))
	assert.EqualError(t, err, "test error 3")
	assert.Equal(t, 0, n)
}

func TestDataChunk_ReadWithDelay(t *testing.T) {
	d := 111 * time.Millisecond
	r := fakereader.Reader{Flow: []io.Reader{
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Delay: d},                                              // 1
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Delay: d},                                              // 2
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true, Delay: d},                                // 3
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true, Delay: d},                                // 4
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true, Delay: d},                                // 5
		&fakereader.DataChunk{Bytes: []byte{1, 2, 3}, Strict: true, Delay: d, Err: errors.New("test error")}, // 6
	}}

	// 1
	start := time.Now()
	n, err := r.Read(make([]byte, 3))
	assert.True(t, time.Since(start) >= d)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	// 2
	start = time.Now()
	n, err = r.Read(make([]byte, 1))
	assert.True(t, time.Since(start) >= d)
	assert.Error(t, err)
	assert.Equal(t, 0, n)

	// 3
	start = time.Now()
	n, err = r.Read(make([]byte, 3))
	assert.True(t, time.Since(start) >= d)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	// 4
	start = time.Now()
	n, err = r.Read(make([]byte, 1))
	assert.True(t, time.Since(start) >= d)
	assert.Error(t, err)
	assert.Equal(t, 0, n)

	// 5
	start = time.Now()
	n, err = r.Read(make([]byte, 7))
	assert.True(t, time.Since(start) >= d)
	assert.Error(t, err)
	assert.Equal(t, 0, n)

	// 6
	start = time.Now()
	n, err = r.Read(make([]byte, 3))
	assert.True(t, time.Since(start) >= d)
	assert.EqualError(t, err, "test error")
	assert.Equal(t, 3, n)
}
