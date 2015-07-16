package utee

import (
	"testing"
	"time"
)

func TestLatch(t *testing.T) {
	latch := NewThrottle(100)

	latch.Acquire()
	latch.Acquire()

	if n := latch.Current(); n != 2 {
		t.Error("current should be 2", n)
	}

	if n := latch.Available(); n != 98 {
		t.Error("available should be 98", n)
	}

	latch.Release()

	if n := latch.Current(); n != 1 {
		t.Error("current should be 1", n)
	}

	if n := latch.Available(); n != 99 {
		t.Error("available should be 99", n)
	}

	latch = NewThrottle(3)

	go func() {
		time.Sleep(2 * time.Second)
		latch.Release()
		latch.Release()
	}()

	latch.Acquire()
	latch.Acquire()
	latch.Acquire()
	latch.Acquire()
	latch.Acquire()
}
