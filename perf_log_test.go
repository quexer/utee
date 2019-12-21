package utee

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewPerfLog(t *testing.T) {
	pl := NewPerfLog(100, logrus.WithField("on", "test"))
	pl.Tick(1)
	time.Sleep(10 * time.Millisecond)
	pl.Done() // show output nothing

	pl = NewPerfLog(20, logrus.WithField("on", "test2"))
	time.Sleep(30 * time.Millisecond)
	pl.Tick("a")
	time.Sleep(1200 * time.Millisecond)
	pl.Tick("b")
	time.Sleep(15 * time.Millisecond)
	pl.Done() // show output log
}
