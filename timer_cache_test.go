package utee

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTimeCache(t *testing.T) {
	tc := NewTimerCache(3, func(k, v interface{}) {
		log.Println("expire", k, v)
	})
	if len := tc.Len(); len != 0 {
		t.Error("len should be 0", len)
	}
	tc.Put(1, 1)
	tc.Put(2, 2)
	time.Sleep(time.Second)
	tc.Put(1, 3)
	if len := tc.Len(); len != 2 {
		t.Error("len should be 2", len)
	}

	if val := tc.Get(1); val != 3 {
		t.Error("1=> should be 3", val)
	}

	cb := func(k, v interface{}) {
		fmt.Println("@k:", k, "@v:", v)
	}

	tc.Loop(cb)

	time.Sleep(4 * time.Second)

	if len := tc.Len(); len != 0 {
		t.Error("len should be 0 after expire", len)
	}
}
