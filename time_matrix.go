package utee

import (
	"log"
	"sync"
)

type tmEntry struct {
	name string
	tick int64
}
type TimeMatrix struct {
	sync.Mutex
	m []*tmEntry
}

func (p *TimeMatrix) Rec(name string) {
	p.Lock()
	p.m = append(p.m, &tmEntry{name, Tick()})
	p.Unlock()
}

func (p *TimeMatrix) Print() {
	p.Lock()
	log.Println("    time matrix")
	if len(p.m) < 2 {
		return
	}
	for i, val := range p.m {
		if i == 0 {
			continue
		}

		log.Printf("    %20v: %8vms\n", p.m[i-1].name+"~"+val.name, val.tick-p.m[i-1].tick)
	}
	p.Unlock()
}
