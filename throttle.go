package utee

import (
	"sync"
)

// Throttle , limit number of concurrent goroutines
type Throttle struct {
	ch chan interface{}
	wg sync.WaitGroup
}

// NewThrottle New Throttle
// max , max number of this throttle
func NewThrottle(max int) *Throttle {
	return &Throttle{
		ch: make(chan interface{}, max),
	}
}

// Acquire , acquire 1
func (p *Throttle) Acquire() {
	p.ch <- nil
	p.wg.Add(1)
}

// Release , release 1
func (p *Throttle) Release() {
	<-p.ch
	p.wg.Done()
}

// Available , available number of this throttle
func (p *Throttle) Available() int {
	return cap(p.ch) - len(p.ch)
}

// Current , current number of this throttle
func (p *Throttle) Current() int {
	return len(p.ch)
}

// Wait , wait until current number become 0
func (p *Throttle) Wait() {
	p.wg.Wait()
}
