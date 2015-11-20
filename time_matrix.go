package utee

import (
	"fmt"
	"log"
	"strings"
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
	if len(p.m) < 2 {
		return
	}

	l := []string{"   time matrix"}
	for i, val := range p.m {
		if i == 0 {
			continue
		}
		s := fmt.Sprintf("%35v: %8vms", p.m[i-1].name+"~"+val.name, val.tick-p.m[i-1].tick)
		l = append(l, s)
	}
	log.Println(strings.Join(l, "\n"))
	p.Unlock()
}
