package fakereader_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/gofaker"
	"github.com/albenik/gofaker/clock"
	"github.com/albenik/gofaker/fakereader"
)

func TestReader_Read_Success(t *testing.T) {
	r := fakereader.New(t,
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
	tt := new(gofaker.FailTriggerTest)
	r := fakereader.New(tt)

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, tt.FailedAsExpected)
	assert.Equal(t, "unexpected 1 read", tt.FailMessage)
}

func TestStrictBytesReader_Success(t *testing.T) {
	r := fakereader.New(t,
		fakereader.StrictBytesReader([]byte{1, 2, 3}),
		fakereader.StrictBytesReader([]byte{1, 2}),
		fakereader.StrictBytesReader([]byte{1}),
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

func TestStrictBytesReader_Fail(t *testing.T) {
	tt := new(gofaker.FailTriggerTest)
	r := fakereader.New(tt,
		fakereader.StrictBytesReader([]byte{211, 212, 213}),
	)

	p := make([]byte, 1)
	n, err := r.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte{211}, p)
	if assert.True(t, tt.FailedAsExpected) {
		assert.Equal(t, "expected buffer length is 3 but actual is 1", tt.FailMessage)
	}
}

func TestDelayRead(t *testing.T) {
	now := time.Date(2018, 1, 1, 12, 0, 0, 0, time.Local)
	fakeclock := clock.NewFakeClock(now).Source()

	r := fakereader.New(t,
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

	r := fakereader.New(t,
		fakereader.DelayRead(3*time.Second, fakereader.StrictBytesReader([]byte{1, 2, 3}), fakeclock),
		fakereader.DelayRead(2*time.Second, fakereader.StrictBytesReader([]byte{1, 2}), fakeclock),
		fakereader.DelayRead(1*time.Second, fakereader.StrictBytesReader([]byte{1}), fakeclock),
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
	tt := new(gofaker.FailTriggerTest)
	now := time.Date(2018, 1, 1, 12, 0, 0, 0, time.Local)
	fakeclock := clock.NewFakeClock(now).Source()

	r := fakereader.New(tt,
		fakereader.DelayRead(3*time.Second, fakereader.StrictBytesReader([]byte{211, 212, 213}), fakeclock),
	)

	p := make([]byte, 1)
	n, err := r.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte{211}, p)
	assert.Equal(t, now.Add(3*time.Second), fakeclock.Now())
	if assert.True(t, tt.FailedAsExpected) {
		assert.Equal(t, "expected buffer length is 3 but actual is 1", tt.FailMessage)
	}
}

func TestReadError(t *testing.T) {
	r := fakereader.New(t,
		fakereader.ReadError(errors.New("custom read error")),
	)
	p := make([]byte, 1)
	n, err := r.Read(p)
	assert.EqualError(t, err, "custom read error")
	assert.Equal(t, 0, n)
}
