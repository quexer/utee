package utee

import "time"

// Tick unix timestamp in millisecond
type Tick int64

// TickToTime convert tick to local time
func (p Tick) ToTime() time.Time {
	tick := int64(p)
	return time.Unix(tick/1e3, (tick%1e3)*1e6)
}

// NewTick create Tick. default value if now at local time
func NewTick(t ...time.Time) Tick {
	if len(t) == 0 {
		return Tick(time.Now().UnixNano() / 1e6)
	}

	return Tick(t[0].UnixNano() / 1e6)
}
