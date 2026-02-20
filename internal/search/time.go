package search

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/timeutil"
)

// ParseFlexibleTime interprets a flexible time string and returns epoch milliseconds.
// Accepts:
//   - Epoch ms directly: "1705363200000"
//   - ISO 8601 / RFC 3339: "2024-01-15T00:00:00Z"
//   - Date only: "2024-01-15" (midnight UTC)
//   - Relative: "-24h", "-7d", "-30m"
//   - "now"
func ParseFlexibleTime(input string) (int64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("empty time value")
	}

	if strings.EqualFold(input, "now") {
		return time.Now().UnixMilli(), nil
	}

	// Relative: starts with "-"
	if strings.HasPrefix(input, "-") {
		dur := input[1:] // strip the "-"
		from, _, err := timeutil.ParseLastEpochMs(dur)
		if err != nil {
			return 0, fmt.Errorf("invalid relative time %q: %w", input, err)
		}
		ms, _ := strconv.ParseInt(from, 10, 64)
		return ms, nil
	}

	// Pure numeric: epoch ms
	if ms, err := strconv.ParseInt(input, 10, 64); err == nil && ms > 1000000000000/10 {
		return ms, nil
	}

	// ISO 8601 / RFC 3339
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t.UnixMilli(), nil
	}

	// Date only: YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", input); err == nil {
		return t.UnixMilli(), nil
	}

	return 0, fmt.Errorf("unrecognized time format %q (use epoch ms, ISO 8601, YYYY-MM-DD, -24h, or now)", input)
}
