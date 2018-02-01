package fakewriter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker"
	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakewriter"
)

func TestWriter_EOF(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "unexpected 1 write", tt.FailMessage)
}

func TestWriter_Locked(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)

	w := fakewriter.New(tt,
		fakewriter.ExpectLen(101),
		fakewriter.ExpectLen(102),
		fakewriter.ExpectLen(103),
	)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data length: 101 expected but 1 recieved", tt.FailMessage)

	n, err = w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "writer locked at 2 write", tt.FailMessage)
}

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

func TestShortWrite_Short_Success(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.ShortWrite(1, fakewriter.ExpectLen(3)),
		fakewriter.ShortWrite(1, fakewriter.ExpectData([]byte{1, 2, 3})),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	n, err = w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestShortWrite_NoShort_Success(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.ShortWrite(3, fakewriter.ExpectLen(2)),
		fakewriter.ShortWrite(3, fakewriter.ExpectData([]byte{1, 2})),
	)

	n, err := w.Write([]byte{2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestShortWrite_FailLen(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)
	w := fakewriter.New(tt,
		fakewriter.ShortWrite(1, fakewriter.ExpectLen(2)),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data length: 2 expected but 3 recieved", tt.FailMessage)
}

func TestShortWrite_FailData(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)
	w := fakewriter.New(tt,
		fakewriter.ShortWrite(1, fakewriter.ExpectData([]byte{1, 2, 3})),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data: [01 02 03] expected but [03 02 01] recieved", tt.FailMessage)
}

func TestDelayWrite_Success(t *testing.T) {
	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New(t,
		fakewriter.DelayWrite(15*time.Second, fakewriter.ExpectLen(3), clk),
		fakewriter.DelayWrite(25*time.Second, fakewriter.ExpectData([]byte{1, 2, 3}), clk),
	)

	n, err := w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, now.Add(15*time.Second), clk.Now())

	n, err = w.Write([]byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, now.Add(40*time.Second), clk.Now())
}

func TestDelayWrite_Fail(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)
	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New(tt,
		fakewriter.DelayWrite(15*time.Second, fakewriter.ExpectData([]byte{1, 2, 3}), clk),
	)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, now.Add(15*time.Second), clk.Now())
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "invalid data: [01 02 03] expected but [01] recieved", tt.FailMessage)
}

func TestWriteError(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.WriteError(errors.New("custom write error")),
	)
	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "custom write error")
	assert.Equal(t, 0, n)
}

func TestForceLen(t *testing.T) {
	w := fakewriter.New(t,
		fakewriter.ForceLen(7, fakewriter.WriteError(errors.New("custom write error"))),
	)
	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "custom write error")
	assert.Equal(t, 7, n)
}
