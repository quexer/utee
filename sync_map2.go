package utee

import "sync"

// SyncMap2 generic sync map
type SyncMap2[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
}

func (p *SyncMap2[K, V]) Put(key K, val V) {
	p.Lock()
	defer p.Unlock()

	if p.m == nil {
		p.m = map[K]V{}
	}
	p.m[key] = val
}

func (p *SyncMap2[K, V]) Remove(key K) {
	p.Lock()
	defer p.Unlock()

	if p.m == nil {
		return
	}
	delete(p.m, key)
}

func (p *SyncMap2[K, V]) Clear() {
	p.Lock()
	defer p.Unlock()

	p.m = nil
}

func (p *SyncMap2[K, V]) Len() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.m)
}

func (p *SyncMap2[K, V]) Get(key K) (V, bool) {
	p.RLock()
	defer p.RUnlock()

	if p.m == nil {
		var v V
		return v, false
	}
	val, ok := p.m[key]
	return val, ok
}

func (p *SyncMap2[K, V]) Keys() []K {
	p.RLock()
	defer p.RUnlock()

	l := make([]K, len(p.m))
	for k := range p.m {
		l = append(l, k)
	}
	return l
}
