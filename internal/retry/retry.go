// Package retry provides exponential backoff utilities for HTTP retries.
package retry

import (
	"math"
	"time"
)

// Wait computes the exponential backoff duration for the given attempt number.
// attempt must be >= 1. The returned duration is clamped between min and max.
//
// The formula is: min * 2^(attempt-1), capped at max.
func Wait(attempt int, min, max time.Duration) time.Duration {
	exp := math.Pow(2, float64(attempt-1))
	d := time.Duration(float64(min) * exp)
	if d > max {
		return max
	}
	return d
}
