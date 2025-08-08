package utee

import (
	"time"
)

// TimeBetween check if check time is between start and end
func TimeBetween(check, start, end time.Time) bool {
	return check.After(start) && check.Before(end)
}
