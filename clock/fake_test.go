package clock_test

import (
	"testing"
	"time"

	"github.com/albenik/gofaker/clock"
	"github.com/stretchr/testify/assert"
)

var testtime = time.Date(2018, 1, 1, 13, 14, 15, 167, time.UTC)

func TestFakeClock_Now(t *testing.T) {
	c := clock.NewFakeClock(testtime).Source()
	assert.Equal(t, testtime, c.Now())
	assert.True(t, c.Now().Equal(testtime))
	assert.False(t, c.Now().Before(testtime))
	assert.False(t, c.Now().After(testtime))
}

func TestFakeClock_Since(t *testing.T) {
	c := clock.NewFakeClock(testtime).Source()
	assert.Equal(t, time.Hour, c.Since(testtime.Add(-time.Hour)))
}

func TestFakeClock_Until(t *testing.T) {
	c := clock.NewFakeClock(testtime).Source()
	assert.Equal(t, time.Hour, c.Until(testtime.Add(time.Hour)))
}

func TestFakeClock_Sleep(t *testing.T) {
	c := clock.NewFakeClock(testtime).Source()
	c.Sleep(time.Second)
	assert.Equal(t, testtime.Add(time.Second), c.Now())
}
