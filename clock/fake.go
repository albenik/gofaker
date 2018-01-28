package clock

import "time"

type FakeClock struct {
	time time.Time
}

func NewFakeClock(start time.Time) *FakeClock {
	return &FakeClock{time: start}
}

func (fake *FakeClock) Source() *Source {
	return &Source{
		Now:   fake.now,
		Since: fake.since,
		Until: fake.until,
		Sleep: fake.sleep,
	}
}

func (fake *FakeClock) Advance(d time.Duration) {
	fake.time = fake.time.Add(d)
}

func (fake *FakeClock) now() time.Time {
	return fake.time
}

func (fake *FakeClock) since(t time.Time) time.Duration {
	return fake.time.Sub(t)
}

func (fake *FakeClock) until(t time.Time) time.Duration {
	return t.Sub(fake.time)
}

func (fake *FakeClock) sleep(d time.Duration) {
	fake.time = fake.time.Add(d)
}
