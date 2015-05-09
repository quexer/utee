package utee

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
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
