package sprite

import (
	"log"
	"strconv"
	"strings"
)

func parseInterval(val interface{}) (int, int, int) {
	switch v := val.(type) {
	case int:
		return v, v, 1
	case float64:
		return int(v), int(v), 1
	case string:
		if n, err := strconv.Atoi(v); err == nil {
			return n, n, 1
		}
		matches := intervalMatcher.FindStringSubmatch(strings.TrimSpace(v))
		if len(matches) != 3 {
			log.Fatalf("Could not parse interval from %s", v)
		}
		low, _ := strconv.Atoi(matches[1])
		high, _ := strconv.Atoi(matches[2])
		if low > high {
			return low, high, -1
		} else {
			return low, high, 1
		}
	default:
		log.Fatalf("Could not parse interval from %v", val)
	}

	return 0, 0, 0
}
