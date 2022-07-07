package utee

import "sync"

// SyncMap , thread safe map
// Deprecated , use SyncMap2 instead
type SyncMap struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (p *SyncMap) Put(key, val interface{}) {
	p.Lock()
	defer p.Unlock()

	if p.m == nil {
		p.m = map[interface{}]interface{}{}
	}
	p.m[key] = val
}

func (p *SyncMap) Remove(key interface{}) {
	p.Lock()
	defer p.Unlock()

	if p.m == nil {
		return
	}
	delete(p.m, key)
}

func (p *SyncMap) Clear() {
	p.Lock()
	defer p.Unlock()

	p.m = nil
}

func (p *SyncMap) Len() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.m)
}

func (p *SyncMap) Get(key interface{}) (interface{}, bool) {
	p.RLock()
	defer p.RUnlock()

	if p.m == nil {
		return nil, false
	}
	val, ok := p.m[key]
	return val, ok
}

func (p *SyncMap) Keys() []interface{} {
	p.RLock()
	defer p.RUnlock()

	l := make([]interface{}, len(p.m))
	for k := range p.m {
		l = append(l, k)
	}
	return l
}
