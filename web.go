package utee

import (
	"encoding/json"
	"net/http"
)

type J map[string]interface{}

/*
 web utilities
*/

type Web struct {
	W http.ResponseWriter
}

func (p *Web) Json(code int, data interface{}) (int, string) {
	b, err := json.Marshal(data)
	Chk(err)
	p.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return code, string(b)
}
func (p *Web) Txt(code int, txt string) (int, string) {
	p.W.Header().Set("Content-Type", "html/text; charset=UTF-8")
	return code, txt
}
