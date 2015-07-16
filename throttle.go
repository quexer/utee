package utee

type Throttle struct {
	ch chan interface{}
}

func NewThrottle(max int) *Throttle {
	return &Throttle{
		ch: make(chan interface{}, max),
	}
}

func (p *Throttle) Acquire() {
	p.ch <- nil
}

func (p *Throttle) Release() {
	<-p.ch
}

func (p *Throttle) Available() int {
	return cap(p.ch) - len(p.ch)
}

func (p *Throttle) Current() int {
	return len(p.ch)
}
