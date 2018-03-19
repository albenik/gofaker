package fakewriter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakewriter"
)

func assertRecover(t *testing.T, v interface{}) {
	r := recover()
	if r == nil {
		t.Fatal("Panic expected")
	}
	assert.Equal(t, v, r)
}

func TestWriter_EOF(t *testing.T) {
	defer assertRecover(t, "test write #1: unexpected [01]")

	w := fakewriter.New("test")

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestWriter_Locked(t *testing.T) {
	defer assertRecover(t, "test write #1: invalid data length 1 [01] (expected 101) @ github.com/albenik/gofaker/fakewriter/writer_test.go:36")

	w := fakewriter.New("test",
		fakewriter.ExpectLen(101),
		fakewriter.ExpectLen(102),
		fakewriter.ExpectLen(103),
	)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
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
	defer assertRecover(t, "test write #1: invalid data length 2 [01 02] (expected 3) @ github.com/albenik/gofaker/fakewriter/writer_test.go:70")

	w := fakewriter.New("test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestAssertLen_MismatchZero(t *testing.T) {
	defer assertRecover(t, "test write #1: invalid data length 0 [] (expected 3) @ github.com/albenik/gofaker/fakewriter/writer_test.go:82")

	w := fakewriter.New("test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestAssertLen_MismatchNil(t *testing.T) {
	defer assertRecover(t, "test write #1: invalid data length 0 [] (expected 3) @ github.com/albenik/gofaker/fakewriter/writer_test.go:94")

	w := fakewriter.New("test",
		fakewriter.ExpectLen(3),
	)

	n, err := w.Write(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
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
	defer assertRecover(t, "test write #1: invalid data [01] (expected [01 02 03]) @ github.com/albenik/gofaker/fakewriter/writer_test.go:126")

	w := fakewriter.New("test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestAssertData_MismatchEmpty(t *testing.T) {
	defer assertRecover(t, "test write #1: invalid data [] (expected [01 02 03]) @ github.com/albenik/gofaker/fakewriter/writer_test.go:138")

	w := fakewriter.New("test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestAssertData_MismatchZero(t *testing.T) {
	defer assertRecover(t, "test write #1: invalid data [] (expected [01 02 03]) @ github.com/albenik/gofaker/fakewriter/writer_test.go:150")

	w := fakewriter.New("test",
		fakewriter.ExpectData([]byte{1, 2, 3}),
	)

	n, err := w.Write(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
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
	defer assertRecover(t, "test write #1: invalid data length 3 [03 02 01] (expected 2) @ github.com/albenik/gofaker/fakewriter/writer_test.go:192")

	w := fakewriter.New("test",
		fakewriter.ShortWrite(1, fakewriter.ExpectLen(2)),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
}

func TestShortWrite_FailData(t *testing.T) {
	defer assertRecover(t, "test write #1: invalid data [03 02 01] (expected [01 02 03]) @ github.com/albenik/gofaker/fakewriter/writer_test.go:204")

	w := fakewriter.New("test",
		fakewriter.ShortWrite(1, fakewriter.ExpectData([]byte{1, 2, 3})),
	)

	n, err := w.Write([]byte{3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
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
	defer assertRecover(t, "test write #1: invalid data [01] (expected [01 02 03]) @ github.com/albenik/gofaker/fakewriter/writer_test.go:239")

	now := time.Now()
	clk := clock.NewFakeClock(now).Source()

	w := fakewriter.New("test",
		fakewriter.DelayWrite(15*time.Second, fakewriter.ExpectData([]byte{1, 2, 3}), clk),
	)

	n, err := w.Write([]byte{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, now.Add(15*time.Second), clk.Now())
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
