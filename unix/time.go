package unix

import "time"

// Time is in milliseconds since epoch. Compatible with Java currentTimeMillis
type Time int64

// Now returns the current time in milliseconds
func Now() Time {
	return From(time.Now())
}

// From returns the time from go time in milliseconds
func From(t time.Time) Time {
	return Time(t.Round(time.Millisecond).UnixNano() / 1e6)
}

// Time returns the go time from unix milliseconds time
func (t Time) Time() time.Time {
	return time.Unix(0, int64(t)*int64(time.Millisecond))
}
