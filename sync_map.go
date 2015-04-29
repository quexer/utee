package utee

import "sync"

type SyncMap struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (p *SyncMap) Put(key, val interface{}) {
	p.Lock()
	defer p.Unlock()

	p.m[key] = val
}

func (p *SyncMap) Len() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.m)
}

func (p *SyncMap) Get(key interface{}) (interface{}, bool) {
	p.RLock()
	defer p.RUnlock()
	val, ok := p.m[key]
	return val, ok
}
