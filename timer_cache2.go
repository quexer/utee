package utee

import (
	"container/heap"
	"sync"
	"time"
)

// An Item is something we manage in a priority queue.
type pqItem2[K comparable, V any] struct {
	key   K
	value V
	ttl   int64 // unix timestamp , in second
	index int
	dead  bool // mark as dead,
}

// A priorityQueue2 implements heap.Interface and holds Items.
type priorityQueue2[K comparable, V any] []*pqItem2[K, V]

func (pq priorityQueue2[K, V]) Len() int { return len(pq) }

func (pq priorityQueue2[K, V]) Less(i, j int) bool {
	return pq[i].ttl < pq[j].ttl
}

func (pq priorityQueue2[K, V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue2[K, V]) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem2[K, V])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue2[K, V]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// TimerCache2 generic timer cache
type TimerCache2[K comparable, V any] struct {
	lock sync.RWMutex

	q   priorityQueue2[K, V]
	m   map[K]*pqItem2[K, V]
	ttl int
}

// NewTimerCache2 create new TimerCache2
// ttl in second
// expireCb,  expire callback
func NewTimerCache2[K comparable, V any](ttl int, expireCb ...func(key K, value V)) *TimerCache2[K, V] {
	tc := &TimerCache2[K, V]{
		q:   priorityQueue2[K, V]{},
		m:   map[K]*pqItem2[K, V]{},
		ttl: ttl,
	}

	var cb func(key K, value V)
	if len(expireCb) > 0 {
		cb = expireCb[0]
	}
	go func() {
		for {
			tc.tryPop(time.Now().Unix(), cb)
			time.Sleep(time.Second)
		}
	}()
	return tc
}

func (p *TimerCache2[K, V]) Put(key K, val V) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	ttl := time.Now().Unix() + int64(p.ttl)
	if old := p.m[key]; old == nil {
		item := &pqItem2[K, V]{
			key:   key,
			value: val,
			ttl:   ttl,
		}
		heap.Push(&p.q, item)
		p.m[key] = item
	} else {
		old.value = val
		old.ttl = ttl
		heap.Fix(&p.q, old.index)
	}
	return true
}

func (p *TimerCache2[K, V]) Get(key K) (V, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if item := p.m[key]; item == nil {
		var out V
		return out, false
	} else {
		return item.value, true
	}
}

// TTL Check ttl (in second)
func (p *TimerCache2[K, V]) TTL(key K) (int64, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if item := p.m[key]; item == nil {
		return 0, false
	} else if item.dead {
		return 0, false
	} else {
		return item.ttl - time.Now().Unix(), true
	}
}

func (p *TimerCache2[K, V]) Remove(key K) (V, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if item, ok := p.m[key]; ok {
		item.dead = true // mark dead
		delete(p.m, key)
		return item.value, true
	} else {
		var out V
		return out, false
	}
}

func (p *TimerCache2[K, V]) Len() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return len(p.m)
}

func (p *TimerCache2[K, V]) tryPop(tick int64, expireCb func(key K, value V)) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for p.q.Len() > 0 {
		item := p.q[0]
		if item.ttl > tick {
			// no expire items
			//			log.Println("no expire items", item.ttl, tick)
			return
		}

		item = heap.Pop(&p.q).(*pqItem2[K, V])
		delete(p.m, item.key)

		// ignore dead item
		if !item.dead && expireCb != nil {
			go expireCb(item.key, item.value)
		}
	}
}

func (p *TimerCache2[K, V]) Keys() []K {
	p.lock.RLock()
	defer p.lock.RUnlock()
	keys := make([]K, 0, len(p.m))
	for k := range p.m {
		keys = append(keys, k)
	}
	return keys
}
