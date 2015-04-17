package utee

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

func MidSlowLog(limit int64) func(*http.Request, martini.Context) {
	return func(req *http.Request, c martini.Context) {
		start := Tick()
		defer func() {
			t := Tick() - start
			if t >= limit {
				log.Printf("[slow] %3vms %s \n", t, req.RequestURI)
			}
		}()
		c.Next()
	}
}
