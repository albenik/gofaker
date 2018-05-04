package fakereader_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakereader"
)

type ft struct {
	msg string
}

func (f *ft) Fatalf(format string, args ...interface{}) {
	f.msg = fmt.Sprintf(format, args...)
}

func TestReader_Read_Success(t *testing.T) {
	r := fakereader.New(t, "test",
		bytes.NewReader([]byte{1, 2, 3}),
		bytes.NewReader([]byte{1, 2}),
		bytes.NewReader([]byte{1}),
	)

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
}

func TestReader_Read_Fail(t *testing.T) {
	f := new(ft)

	r := fakereader.New(f, "test")

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, "test read #1: unexpected", f.msg)
}

func TestStrictBytesReader_Success(t *testing.T) {
	r := fakereader.New(t, "test",
		fakereader.NotLessData([]byte{1, 2, 3}),
		fakereader.NotLessData([]byte{1, 2, 3}),
		fakereader.NotLessData([]byte{1, 2}),
		fakereader.NotLessData([]byte{1}),
	)

	readExpecting := func(i int, exp []byte) {
		buf := make([]byte, len(exp))
		n, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, i, n)
		assert.Equal(t, exp, buf)
	}

	readExpecting(3, []byte{1, 2, 3, 0})
	readExpecting(3, []byte{1, 2, 3})
	readExpecting(2, []byte{1, 2})
	readExpecting(1, []byte{1})
}

func TestStrictBytesReader_Fail(t *testing.T) {
	f := new(ft)

	r := fakereader.New(f, "test",
		fakereader.NotLessData([]byte{211, 212, 213}),
	)

	p := []byte{0xFF}
	n, err := r.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, []byte{0xFF}, p)
	assert.NoError(t, err)
	assert.Equal(t, "test read #1 @ github.com/albenik/gofaker/fakereader/reader_test.go:82: wrong destination size 1 (3 expected)", f.msg)
}

func TestUltraStrictBytesReader_Fail(t *testing.T) {
	f := new(ft)

	r := fakereader.New(f, "test",
		fakereader.EqualData([]byte{211, 212, 213}),
	)
	p := []byte{0xFF}
	n, err := r.Read(p)

	assert.Equal(t, 0, n)
	assert.Equal(t, []byte{0xFF}, p)
	assert.NoError(t, err)
	assert.Equal(t, "test read #1 @ github.com/albenik/gofaker/fakereader/reader_test.go:97: wrong destination size 1 (3 expected)", f.msg)
}

func TestUltraStrictBytesReader_Fail2(t *testing.T) {
	f := new(ft)

	r := fakereader.New(f, "test",
		fakereader.EqualData([]byte{211, 212, 213}),
	)
	p := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	n, err := r.Read(p)

	assert.Equal(t, 0, n)
	assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF}, p)
	assert.NoError(t, err)
	assert.Equal(t, "test read #1 @ github.com/albenik/gofaker/fakereader/reader_test.go:112: wrong destination size 4 (3 expected)", f.msg)
}

func TestDelayRead(t *testing.T) {
	now := time.Date(2018, 1, 1, 12, 0, 0, 0, time.Local)
	fakeclock := clock.NewFakeClock(now).Source()

	r := fakereader.New(t, "test",
		fakereader.DelayRead(3*time.Second, bytes.NewReader([]byte{1, 2, 3}), fakeclock),
		fakereader.DelayRead(2*time.Second, bytes.NewReader([]byte{1, 2}), fakeclock),
		fakereader.DelayRead(1*time.Second, bytes.NewReader([]byte{1}), fakeclock),
	)

	read := func(expectData []byte, expectTime time.Time) {
		buf := make([]byte, len(expectData))
		n, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, len(expectData), n)
		assert.Equal(t, expectData, buf)
		assert.Equal(t, expectTime, fakeclock.Now())
	}

	read([]byte{1, 2, 3}, now.Add(3*time.Second))
	read([]byte{1, 2}, now.Add(5*time.Second))
	read([]byte{1}, now.Add(6*time.Second))
}

func TestDelayRead_Strict_Success(t *testing.T) {
	now := time.Date(2018, 1, 1, 12, 0, 0, 0, time.Local)
	fakeclock := clock.NewFakeClock(now).Source()

	r := fakereader.New(t, "test",
		fakereader.DelayRead(3*time.Second, fakereader.NotLessData([]byte{1, 2, 3}), fakeclock),
		fakereader.DelayRead(2*time.Second, fakereader.NotLessData([]byte{1, 2}), fakeclock),
		fakereader.DelayRead(1*time.Second, fakereader.NotLessData([]byte{1}), fakeclock),
	)

	read := func(expectData []byte, expectTime time.Time) {
		buf := make([]byte, len(expectData))
		n, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, len(expectData), n)
		assert.Equal(t, expectData, buf)
		assert.Equal(t, expectTime, fakeclock.Now())
	}

	read([]byte{1, 2, 3}, now.Add(3*time.Second))
	read([]byte{1, 2}, now.Add(5*time.Second))
	read([]byte{1}, now.Add(6*time.Second))
}

func TestDelayRead_Strict_Fail(t *testing.T) {
	now := time.Date(2018, 1, 1, 12, 0, 0, 0, time.Local)
	fakeclock := clock.NewFakeClock(now).Source()

	f := new(ft)

	r := fakereader.New(f, "test",
		fakereader.DelayRead(3*time.Second, fakereader.NotLessData([]byte{211, 212, 213}), fakeclock),
	)

	p := make([]byte, 1)
	n, err := r.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, now.Add(3*time.Second), fakeclock.Now())
	assert.Equal(t, 0, n)
	assert.Equal(t, []byte{0x00}, p)
	assert.Equal(t, "test read #1 @ github.com/albenik/gofaker/fakereader/reader_test.go:178: wrong destination size 1 (3 expected)", f.msg)
}

func TestReadError(t *testing.T) {
	r := fakereader.New(t, "test",
		fakereader.ReturnError(errors.New("custom read error")),
	)
	p := make([]byte, 1)
	n, err := r.Read(p)
	assert.Equal(t, 0, n)
	assert.EqualError(t, err, "custom read error")
}
