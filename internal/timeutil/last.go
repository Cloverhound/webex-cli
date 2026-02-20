package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseDuration parses a shorthand like "5m", "1h", "24h", "7d" into a time.Duration.
func parseDuration(last string) (time.Duration, error) {
	last = strings.TrimSpace(last)
	if len(last) < 2 {
		return 0, fmt.Errorf("invalid --last value %q: expected number + unit (e.g. 1h, 30m)", last)
	}

	unit := last[len(last)-1]
	numStr := last[:len(last)-1]

	n, err := strconv.ParseFloat(numStr, 64)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid --last value %q: expected positive number + unit (e.g. 1h, 30m)", last)
	}

	switch unit {
	case 'm':
		return time.Duration(n * float64(time.Minute)), nil
	case 'h':
		return time.Duration(n * float64(time.Hour)), nil
	case 'd':
		return time.Duration(n * 24 * float64(time.Hour)), nil
	default:
		return 0, fmt.Errorf("invalid --last unit %q: expected m (minutes), h (hours), or d (days)", string(unit))
	}
}

// ParseLastEpochMs parses a duration shorthand and returns epoch millisecond
// strings for from (now - duration) and to (now).
func ParseLastEpochMs(last string) (string, string, error) {
	d, err := parseDuration(last)
	if err != nil {
		return "", "", err
	}
	now := time.Now()
	return strconv.FormatInt(now.Add(-d).UnixMilli(), 10), strconv.FormatInt(now.UnixMilli(), 10), nil
}

// ParseLastISO parses a duration shorthand and returns ISO 8601 (RFC 3339)
// strings for from (now - duration) and to (now).
func ParseLastISO(last string) (string, string, error) {
	d, err := parseDuration(last)
	if err != nil {
		return "", "", err
	}
	now := time.Now().UTC()
	return now.Add(-d).Format(time.RFC3339), now.Format(time.RFC3339), nil
}
