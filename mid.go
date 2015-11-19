package utee

import (
	"expvar"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
	"time"
)

func MidSlowLog(limit int) func(*http.Request, martini.Context) {
	if limit <= 0 {
		log.Fatalln("[slow log] err:  bad limit")
	}

	return func(req *http.Request, c martini.Context) {
		start := Tick()
		tm := &TimeMatrix{}
		tm.Rec("start")
		c.Map(tm)
		defer func() {
			t := Tick() - start
			if t >= int64(limit) {
				log.Printf("[slow] %3vms %s \n", t, req.RequestURI)
				if Env("TIME_MATRIX", false, false) != "" {
					tm.Print()
				}
			}
		}()
		c.Next()
		tm.Rec("end")
	}
}

func MidWeb(w http.ResponseWriter, c martini.Context) {
	web := &Web{W: w}
	c.Map(web)
}

func MidTextDefault(w http.ResponseWriter) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	}
}

var (
	expServerConcurrent = expvar.NewInt("z_utee_server_concurrent")
	expServeCount       = expvar.NewInt("z_utee_serve_count")
	expTps              = expvar.NewInt("z_utee_serve_tps")
)

func MidConcurrent(concurrent ...int) func(http.ResponseWriter, martini.Context) {
	go func() {
		lastSecond, _ := strconv.ParseInt(expServeCount.String(), 10, 64)
		for range time.Tick(time.Second) {
			countTotal, _ := strconv.ParseInt(expServeCount.String(), 10, 64)
			expTps.Set(countTotal - lastSecond)
			lastSecond = countTotal
		}
	}()

	var ch chan string
	if len(concurrent) > 0 {
		if concurrent[0] > 0 {
			ch = make(chan string, concurrent[0])
		} else {
			log.Fatalln("bad concurrent number", concurrent[0])
		}
	}
	return func(w http.ResponseWriter, c martini.Context) {
		expServerConcurrent.Add(1)
		defer func() {
			expServerConcurrent.Add(-1)
			expServeCount.Add(1)
		}()
		if ch == nil {
			c.Next()
		} else {
			select {
			case ch <- "a":
				func() {
					defer func() {
						<-ch
					}()
					c.Next()
				}()
			default:
				log.Println("[warn] reach concurrent limit:", concurrent[0])
				http.Error(w, "server is busy", 503)
			}
		}
	}
}
