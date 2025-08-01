package utee

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// TimeBetween check if check time is between start and end
func TimeBetween(check, start, end time.Time) bool {
	return check.After(start) && check.Before(end)
}

// Step records information for each step
type Step struct {
	Name     string        // step name
	Duration time.Duration // duration
}

// TimeTracker records the duration of a series of steps
type TimeTracker struct {
	sync.Mutex
	startTime time.Time // start time
	lastTime  time.Time // last recorded time
	steps     []Step    // list of steps
}

// NewTimeTracker creates a new TimeTracker instance
func NewTimeTracker() *TimeTracker {
	now := time.Now()
	return &TimeTracker{
		startTime: now,
		lastTime:  now,
		steps:     make([]Step, 0, 10), // initial capacity is 10
	}
}

// Tick records the current step's name and duration
func (p *TimeTracker) Tick(stepName string) {
	p.Lock()
	defer p.Unlock()

	now := time.Now()
	duration := now.Sub(p.lastTime)

	p.steps = append(p.steps, Step{
		Name:     stepName,
		Duration: duration,
	})

	p.lastTime = now
}

// Total returns the total duration
func (p *TimeTracker) Total() time.Duration {
	return p.lastTime.Sub(p.startTime)
}

// Steps returns the list of all steps for custom formatting
func (p *TimeTracker) Steps() []Step {
	return Clone(p.steps)
}

// ToString returns the string representation of all steps
// The format is "|-duration-stepName-duration-stepName"
func (p *TimeTracker) ToString() string {
	p.Lock()
	defer p.Unlock()

	if len(p.steps) == 0 {
		return "|"
	}

	var builder strings.Builder
	builder.WriteString("|")

	for _, v := range p.steps {
		// Round the time, removing the part below milliseconds
		duration := v.Duration.Truncate(time.Millisecond)
		builder.WriteString(fmt.Sprintf("-%v-%s", duration, v.Name))
	}

	return builder.String()
}
