package balancer

import "time"

// IsWithin accepts a Unix seconds-based timestamp for start and end, and
// returns true if the provided durationFromStart has elapsed since end.
// End is typically time.Now().Unix().
func IsWithin(start, end int64, durationFromStart time.Duration) bool {
	return time.Unix(start, 0).Add(durationFromStart).After(time.Unix(end, 0))
}
