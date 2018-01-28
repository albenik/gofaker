package fakewriter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker"
	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakewriter"
)

func TestAssertLen_OK(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.ExpectLen(3),
		fakewriter.ExpectLen(2),
		fakewriter.ExpectLen(1),
	)

	n, err := w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestAssertLen_MismatchNonzero(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data length: 3 expected but 2 recieved", tt.FailMessage)
}

func TestAssertLen_MismatchZero(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data length: 3 expected but 0 recieved", tt.FailMessage)
}

func TestAssertLen_MismatchNil(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data length: 3 expected but 0 recieved", tt.FailMessage)
}

func TestAssertData_OK(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.ExpectData([]byte{1, 2, 3}),
		fakewriter.ExpectData([]byte{1, 2}),
		fakewriter.ExpectData([]byte{1}),
	)

	n, err := w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestAssertData_MismatchNonempty(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data: [01 02 03] expected but [01] recieved", tt.FailMessage)
}

func TestAssertData_MismatchEmpty(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data: [01 02 03] expected but [] recieved", tt.FailMessage)
}

func TestAssertData_MismatchZero(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data: [01 02 03] expected but [] recieved", tt.FailMessage)
}

func TestWriter_UnexpectedExtraWrite(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt, nil)

	n, err := w.Write([]byte{1})

	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "unexpected 1 write", tt.FailMessage)
}

func TestTruncateWrite(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.ShortWrite(3),
		fakewriter.ShortWrite(3),
	)

	n, err := w.Write([]byte{1, 2, 3, 4, 5, 6, 7})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)

	n, err = w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestDelayWrite(t *testing.T) {
	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New(t,
		fakewriter.DelayWrite(15*time.Second, clk),
		fakewriter.DelayWrite(25*time.Second, clk),
	)

	n, err := w.Write([]byte{1, 2, 3, 4, 5, 6, 7})
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, now.Add(15*time.Second), clk.Now())

	n, err = w.Write([]byte{1, 2, 3, 4, 5, 6, 7})
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, now.Add(40*time.Second), clk.Now())
}
