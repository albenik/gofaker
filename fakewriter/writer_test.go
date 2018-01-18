package fakewriter_test

import (
	"io"
	"testing"

	"github.com/albenik/gofaker/fakewriter"
	"github.com/stretchr/testify/assert"
)

func TestAssertLen(t *testing.T) {
	w := fakewriter.Writer{Flow: []io.Writer{
		fakewriter.AssertLen(3),
		fakewriter.AssertLen(2),
		fakewriter.AssertLen(1),
		fakewriter.AssertLen(2),
		fakewriter.AssertLen(1),
		fakewriter.AssertLen(1),
	}}

	n, err := w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	n, err = w.Write([]byte{1})
	assert.EqualError(t, err, "invalid data length: 2 expected but 1 recieved")
	assert.Equal(t, 0, n)

	n, err = w.Write([]byte{})
	assert.EqualError(t, err, "invalid data length: 1 expected but 0 recieved")
	assert.Equal(t, 0, n)

	n, err = w.Write(nil)
	assert.EqualError(t, err, "invalid data length: 1 expected but 0 recieved")
	assert.Equal(t, 0, n)
}

func TestAssertData(t *testing.T) {
	w := fakewriter.Writer{Flow: []io.Writer{
		fakewriter.AssertData([]byte{1, 2, 3}),
		fakewriter.AssertData([]byte{1, 2}),
		fakewriter.AssertData([]byte{1}),
		fakewriter.AssertData([]byte{1, 2}),
		fakewriter.AssertData([]byte{1}),
		fakewriter.AssertData([]byte{1}),
	}}

	n, err := w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	n, err = w.Write([]byte{1})
	assert.EqualError(t, err, "invalid data: [01 02] expected but [01] recieved")
	assert.Equal(t, 0, n)

	n, err = w.Write([]byte{})
	assert.EqualError(t, err, "invalid data: [01] expected but [] recieved")
	assert.Equal(t, 0, n)

	n, err = w.Write(nil)
	assert.EqualError(t, err, "invalid data: [01] expected but [] recieved")
	assert.Equal(t, 0, n)
}

func TestWriter_Write(t *testing.T) {
	w := fakewriter.Writer{}
	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "unexpected 1 write")
	assert.Equal(t, 0, n)
}
