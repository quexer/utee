package utee

import (
	"errors"
	"log"
)

var ErrFull = errors.New("queue is full")

type MemQueue chan interface{}

func NewMemQueue(cap int) MemQueue {
	return make(chan interface{}, cap)
}

//create memory queue, auto-leak element concurrently to worker
func NewLeakMemQueue(cap, concurrent int, worker func(interface{})) MemQueue {
	q := NewMemQueue(cap)

	f := func() {
		for {
			worker(q.Deq())
		}
	}

	for i := 0; i < concurrent; i++ {
		go f()
	}
	return q
}

//enqueue, block if queue is full
func (p MemQueue) EnqBlocking(data interface{}) {
	p <- data
}

//enqueue, return error if queue is full
func (p MemQueue) Enq(data interface{}) error {
	select {
	case p <- data:
	default:
		return ErrFull
	}
	return nil
}

func (p MemQueue) Deq() interface{} {
	return <-p
}

//dequeue less than n in a batch
func (p MemQueue) DeqN(n int) []interface{} {
	if n <= 0 {
		log.Println("[MemQueue] deqn err, n must > 0")
		return nil
	}

	var l []interface{}

	for {
		select {
		case data := <-p:
			l = append(l, data)
			if len(l) == n {
				return l
			}
		default:
			return l
		}
	}
}

func (p MemQueue) Len() int {
	return len(p)
}

func (p MemQueue) Cap() int {
	return cap(p)
}
