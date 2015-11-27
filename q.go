package utee

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
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

	select {
	case data := <-p:
		l = append(l, data)
		if len(l) == n {
			break
		}
	default:
		break
	}

	return l
}

func (p MemQueue) Len() int {
	return len(p)
}

func (p MemQueue) Cap() int {
	return cap(p)
}

func qname(name string) string {
	return "q" + name
}

type SimpleRedisQueue struct {
	name   string
	pool   *redis.Pool
	buffer MemQueue
	batch  int
}

//redis queue with optional memory buffer
//server: redis server address
//auth: redis auth
//name: queue name in redis
//concurrent: concurrent number redis enqueue operation, must >=1
//enqBatch: batch enqueue number, must >=1
//buffer: memory buffer capacity, must >= 0
func NewSimpleRedisQueue(server, auth, name string, concurrent, enqBatch, buffer int) *SimpleRedisQueue {
	if concurrent < 1 {
		log.Fatal("concurrent must >= 1")
	}
	if enqBatch < 1 {
		log.Fatal("batch must >= 1")
	}

	if buffer < 0 {
		log.Fatal("buffer must >= 0")
	}

	q := &SimpleRedisQueue{
		name:   qname(name),
		pool:   CreateRedisPool(concurrent, server, auth),
		buffer: NewMemQueue(buffer),
		batch:  enqBatch,
	}
	for i := 0; i < concurrent; i++ {
		go q.enqLoop()
	}
	return q
}

func (p *SimpleRedisQueue) enqLoop() {
	for {
		l := p.buffer.DeqN(p.batch)
		if len(l) > 0 {
			if err := p.enqBatch(l); err != nil {
				log.Println("[SimpleRedisQueue enqLoop] err ", err)
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *SimpleRedisQueue) enqBatch(l []interface{}) error {
	c := p.pool.Get()
	defer c.Close()
	for _, data := range l {
		if err := c.Send("RPUSH", p.name, data); err != nil {
			log.Println("[SimpleRedisQueue enqBatch] err :", err)
		}
	}
	return c.Flush()
}

func (p *SimpleRedisQueue) Len() (int, error) {
	c := p.pool.Get()
	defer c.Close()

	i, err := redis.Int(c.Do("LLEN", p.name))

	if err != nil && err.Error() == "redigo: nil returned" {
		//expire
		return 0, nil
	}
	return i, err
}

//enqueue, block if buffer is full
func (p *SimpleRedisQueue) EnqBlocking(data interface{}) {
	p.buffer.EnqBlocking(data)
}

//enqueue, return error if buffer is full
func (p *SimpleRedisQueue) Enq(data interface{}) error {
	return p.buffer.Enq(data)
}

func (p *SimpleRedisQueue) Deq() (interface{}, error) {
	c := p.pool.Get()
	defer c.Close()
	return c.Do("LPOP", p.name)
}

func (p *SimpleRedisQueue) BufferLen() int {
	return p.buffer.Len()
}

func (p *SimpleRedisQueue) BufferCap() int {
	return p.buffer.Cap()
}

func (p *SimpleRedisQueue) PollSize() int {
	return p.pool.ActiveCount()
}
