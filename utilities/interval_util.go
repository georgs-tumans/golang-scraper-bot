package utilities

import "time"

// Parses durations from a string including support for "d" as days.
//
// Available time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d".
//
// Will default to 24 hours when failing to parse a duration ending with "d".
func ParseDurationWithDays(s string) (time.Duration, error) {
	if len(s) > 0 && s[len(s)-1] == 'd' {
		// If the duration ends with 'd', parse it as days
		days, err := time.ParseDuration(s[:len(s)-1] + "h")
		if err != nil {
			return 0, err
		}
		return days * 24, nil
	}

	return time.ParseDuration(s)
}
