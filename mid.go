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
		defer func() {
			t := Tick() - start
			if t >= int64(limit) {
				log.Printf("[slow] %3vms %s \n", t, req.RequestURI)
			}
		}()
		c.Next()
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

func MidConcurrent() func(martini.Context) {
	go func() {
		lastSecond, _ := strconv.ParseInt(expServeCount.String(), 10, 64)
		for range time.Tick(time.Second) {
			countTotal, _ := strconv.ParseInt(expServeCount.String(), 10, 64)
			expTps.Set(countTotal - lastSecond)
			lastSecond = countTotal
		}
	}()

	return func(c martini.Context) {
		expServerConcurrent.Add(1)
		defer func() {
			expServerConcurrent.Add(-1)
			expServeCount.Add(1)
		}()
		c.Next()
	}
}
