package clock

import "time"

type Source struct {
	Now   func() time.Time
	Since func(time.Time) time.Duration
	Until func(time.Time) time.Duration
	Sleep func(time.Duration)
}

var Native = Source{
	Now:   time.Now,
	Since: time.Since,
	Sleep: time.Sleep,
	Until: time.Until,
}
