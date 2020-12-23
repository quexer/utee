package utee

import (
	"container/heap"
	"sync"
	"time"
)

// An Item is something we manage in a priority queue.
type pqItem struct {
	key   interface{}
	value interface{}
	ttl   int64 // unix timestamp , in second
	index int
	dead  bool // mark as dead,
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].ttl < pq[j].ttl
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type TimerCache struct {
	lock sync.RWMutex

	q   priorityQueue
	m   map[interface{}]*pqItem
	ttl int
}

// ttl in second
// expireCb,  expire callback
func NewTimerCache(ttl int, expireCb func(key, value interface{})) *TimerCache {
	tc := &TimerCache{
		q:   []*pqItem{},
		m:   map[interface{}]*pqItem{},
		ttl: ttl,
	}

	go func() {
		for {
			tc.tryPop(time.Now().Unix(), expireCb)
			time.Sleep(time.Second)
		}
	}()
	return tc
}

func (p *TimerCache) Put(key, val interface{}) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	ttl := time.Now().Unix() + int64(p.ttl)
	if old := p.m[key]; old == nil {
		item := &pqItem{
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

func (p *TimerCache) Get(key interface{}) interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if item := p.m[key]; item == nil {
		return nil
	} else {
		return item.value
	}
}

// TTL Check ttl (in second)
func (p *TimerCache) TTL(key interface{}) int64 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if item := p.m[key]; item == nil {
		return 0
	} else if item.dead {
		return 0
	} else {
		return item.ttl - time.Now().Unix()
	}
}

func (p *TimerCache) Remove(key interface{}) interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()

	if item, ok := p.m[key]; ok {
		item.dead = true // mark dead
		delete(p.m, key)
		return item.value
	} else {
		return nil
	}
}

func (p *TimerCache) Len() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return len(p.m)
}

func (p *TimerCache) tryPop(tick int64, expireCb func(key, value interface{})) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for p.q.Len() > 0 {
		item := p.q[0]
		if item.ttl > tick {
			// no expire items
			//			log.Println("no expire items", item.ttl, tick)
			return
		}

		item = heap.Pop(&p.q).(*pqItem)
		delete(p.m, item.key)

		// ignore dead item
		if !item.dead && expireCb != nil {
			go expireCb(item.key, item.value)
		}
	}
}

func (p *TimerCache) Keys() []interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()
	keys := make([]interface{}, 0, len(p.m))
	for k := range p.m {
		keys = append(keys, k)
	}
	return keys
}
