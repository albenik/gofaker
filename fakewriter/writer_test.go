package fakewriter_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakewriter"
)

type ft struct {
	msg string
}

func (f *ft) Fatalf(format string, args ...interface{}) {
	f.msg = fmt.Sprintf(format, args...)
}

func TestWriter_EOF(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test")

	n, err := w.Write([]byte{1})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1: unexpected [01]", f.msg)
}

func TestWriter_Locked(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectLen(101),
		fakewriter.ExpectLen(102),
		fakewriter.ExpectLen(103),
	)

	n, err := w.Write([]byte{1})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:38: invalid data length 1 [01] (expected 101)", f.msg)
}

func TestAssertLen_OK(t *testing.T) {
	w := fakewriter.New(t, "test",
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
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{1, 2})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:73: invalid data length 2 [01 02] (expected 3)", f.msg)
}

func TestAssertLen_MismatchZero(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:86: invalid data length 0 [] (expected 3)", f.msg)
}

func TestAssertLen_MismatchNil(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write(nil)
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:99: invalid data length 0 [] (expected 3)", f.msg)
}

func TestAssertData_OK(t *testing.T) {
	w := fakewriter.New(t, "test",
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
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{1})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:132: invalid data [01] (expected [01 02 03])", f.msg)
}

func TestAssertData_MismatchEmpty(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:145: invalid data [] (expected [01 02 03])", f.msg)
}

func TestAssertData_MismatchZero(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write(nil)
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:158: invalid data [] (expected [01 02 03])", f.msg)
}

func TestShortWrite_Short_Success(t *testing.T) {
	w := fakewriter.New(t, "test",
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
	w := fakewriter.New(t, "test",
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
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ShortWrite(1, fakewriter.ExpectLen(2)),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:201: invalid data length 3 [03 02 01] (expected 2)", f.msg)
}

func TestShortWrite_FailData(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ShortWrite(1, fakewriter.ExpectData([]byte{1, 2, 3})),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:214: invalid data [03 02 01] (expected [01 02 03])", f.msg)
}

func TestDelayWrite_Success(t *testing.T) {
	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New(t, "test",
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

	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.DelayWrite(15*time.Second, fakewriter.ExpectData([]byte{1, 2, 3}), clk),
	)

	n, err := w.Write([]byte{1})
	assert.Equal(t, now.Add(15*time.Second), clk.Now())
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test write #1 @ github.com/albenik/gofaker/fakewriter/writer_test.go:250: invalid data [01] (expected [01 02 03])", f.msg)
}

func TestWriteError(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.WriteError(errors.New("custom write error")),
	)

	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "custom write error")
	assert.Equal(t, 0, n)
}

func TestForceLen(t *testing.T) {
	f := new(ft)

	w := fakewriter.New(f, "test",
		fakewriter.ForceLen(7, fakewriter.WriteError(errors.New("custom write error"))),
	)

	n, err := w.Write([]byte{1})
	assert.EqualError(t, err, "custom write error")
	assert.Equal(t, 7, n)
}
