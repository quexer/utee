package utee

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type tick struct {
	Label interface{}
	T     time.Time
}

//性能日志
type PerfLog struct {
	sync.Mutex
	logger *logrus.Entry
	max    uint32  //阈值， 以毫秒为单位，总时间超过此数值才会有日志输出
	ticks  []*tick //记录所有时间点及label
	done   bool
}

func (p *PerfLog) Tick(label interface{}) {
	p.Lock()
	defer p.Unlock()
	p.saveTick(label)
}

func (p *PerfLog) Done() {
	p.Lock()
	defer p.Unlock()

	if p.done {
		p.logger.Errorln("perf_log_already_done")
		return
	}
	p.done = true
	p.saveTick("done")

	//有始有终， 所以至少会记录两个点
	start := p.ticks[0].T
	end := p.ticks[len(p.ticks)-1].T

	totalElapsed := end.Sub(start)
	if totalElapsed < time.Duration(p.max)*time.Millisecond {
		return
	}

	sb := strings.Builder{}
	for i, v := range p.ticks {
		sb.WriteString(fmt.Sprint(`'`, v.Label, `'`))
		if i == len(p.ticks)-1 {
			//最后一轮，不输出时间差
			continue
		}
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(int(p.ticks[i+1].T.Sub(v.T) / time.Millisecond)))
		sb.WriteString("ms ")
	}

	p.logger.WithField("perf_log_total", totalElapsed.String()).Warnln("perf_log_slow", sb.String())
}

//生成PerfLog， maxMs 输出阈值，单位为毫秒
func NewPerfLog(maxMs uint32, logger *logrus.Entry) *PerfLog {
	if maxMs == 0 {
		panic("bad maxMs val")
	}

	pl := &PerfLog{
		logger: logger,
		max:    maxMs,
	}

	pl.Tick("start")
	return pl
}

func (p *PerfLog) saveTick(label interface{}) {
	p.ticks = append(p.ticks, &tick{Label: label, T: time.Now()})
}
