package fakewriter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/albenik/gofaker"
	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakewriter"
)

func TestWriter_EOF(t *testing.T) {
	w := fakewriter.New("test")

	n, err := w.Write([]byte{1})
	assert.Equal(t, 0, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1: unexpected [01]",
	}, err)
}

func TestWriter_Locked(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ExpectLen(101),
		fakewriter.ExpectLen(102),
		fakewriter.ExpectLen(103),
	)

	n, err := w.Write([]byte{1})
	assert.Equal(t, 1, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:27: invalid data length 1 [01] (expected 101)",
	}, err)
}

func TestAssertLen_OK(t *testing.T) {
	w := fakewriter.New("test",
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
	w := fakewriter.New("test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{1, 2})
	assert.Equal(t, 2, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:61: invalid data length 2 [01 02] (expected 3)",
	}, err)
}

func TestAssertLen_MismatchZero(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{})
	assert.Equal(t, 0, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:73: invalid data length 0 [] (expected 3)",
	}, err)
}

func TestAssertLen_MismatchNil(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write(nil)
	assert.Equal(t, 0, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:85: invalid data length 0 [] (expected 3)",
	}, err)
}

func TestAssertData_OK(t *testing.T) {
	w := fakewriter.New("test",
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
	w := fakewriter.New("test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{1})
	assert.Equal(t, 1, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:117: invalid data [01] (expected [01 02 03])",
	}, err)
}

func TestAssertData_MismatchEmpty(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{})
	assert.Equal(t, 0, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:129: invalid data [] (expected [01 02 03])",
	}, err)
}

func TestAssertData_MismatchZero(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write(nil)
	assert.Equal(t, 0, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:141: invalid data [] (expected [01 02 03])",
	}, err)
}

func TestShortWrite_Short_Success(t *testing.T) {
	w := fakewriter.New("test",
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
	w := fakewriter.New("test",
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
	w := fakewriter.New("test",
		fakewriter.ShortWrite(1, fakewriter.ExpectLen(2)),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.Equal(t, 3, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:183: invalid data length 3 [03 02 01] (expected 2)",
	}, err)
}

func TestShortWrite_FailData(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ShortWrite(1, fakewriter.ExpectData([]byte{1, 2, 3})),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.Equal(t, 3, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:195: invalid data [03 02 01] (expected [01 02 03])",
	}, err)
}

func TestDelayWrite_Success(t *testing.T) {
	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New("test",
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
	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New("test",
		fakewriter.DelayWrite(15*time.Second, fakewriter.ExpectData([]byte{1, 2, 3}), clk),
	)

	n, err := w.Write([]byte{1})
	assert.Equal(t, now.Add(15*time.Second), clk.Now())
	assert.Equal(t, 1, n)
	assert.Equal(t, &gofaker.AssertionFailedError{
		Message: "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:230: invalid data [01] (expected [01 02 03])",
	}, err)
}

func TestWriteError(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.WriteError(errors.New("custom write error")),
	)
	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "custom write error")
	assert.Equal(t, 0, n)
}

func TestForceLen(t *testing.T) {
	w := fakewriter.New("test",
		fakewriter.ForceLen(7, fakewriter.WriteError(errors.New("custom write error"))),
	)
	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "custom write error")
	assert.Equal(t, 7, n)
}
