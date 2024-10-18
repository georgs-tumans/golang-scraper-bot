package utilities

import (
	"errors"
	"regexp"
	"time"
)

// Parses durations from a string including support for "d" as days.
//
// Available time units are "m", "h", "d". Anything other than these is not allowed.
//
// Will default to 24 hours when failing to parse a duration ending with "d".
func ParseDurationWithDays(s string) (time.Duration, error) {
	intervalType := extractLetter(s)
	if intervalType != "m" && intervalType != "h" && intervalType != "d" {
		return 0, errors.New("invalid_interval_value")
	}

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

func extractLetter(s string) string {
	// Matches one or more letters
	re := regexp.MustCompile("[a-zA-Z]+")

	return re.FindString(s)
}
