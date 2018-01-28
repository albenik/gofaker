package fakereader_test

import (
	"bytes"
	"fmt"
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
		fakereader.StrictBytesReader([]byte{1, 2, 3}),
		fakereader.StrictBytesReader([]byte{1, 2}),
		fakereader.StrictBytesReader([]byte{1}),
	)

	read := func(expected, actual int) {
		buf := make([]byte, expected)
		n, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, actual, n)
		if assert.True(t, tt.FailedAsExpected) {
			assert.Equal(t, fmt.Sprintf("expected buffer length is %d but actual is %d", actual, expected), tt.FailMessage)
		}
	}

	read(7, 3)
	read(6, 2)
	read(5, 1)
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
		fakereader.DelayRead(3*time.Second, fakereader.StrictBytesReader([]byte{1, 2, 3}), fakeclock),
		fakereader.DelayRead(2*time.Second, fakereader.StrictBytesReader([]byte{1, 2}), fakeclock),
		fakereader.DelayRead(1*time.Second, fakereader.StrictBytesReader([]byte{1}), fakeclock),
	)

	read := func(expectedLen, actualLen int, expectTime time.Time) {
		buf := make([]byte, expectedLen)
		n, err := r.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, actualLen, n)
		assert.Equal(t, expectTime, fakeclock.Now())
		if assert.True(t, tt.FailedAsExpected) {
			assert.Equal(t, fmt.Sprintf("expected buffer length is %d but actual is %d", actualLen, expectedLen), tt.FailMessage)
		}
	}

	read(7, 3, now.Add(3*time.Second))
	read(6, 2, now.Add(5*time.Second))
	read(5, 1, now.Add(6*time.Second))
}
