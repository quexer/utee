package utee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

type HttpTestSuite struct {
	Name     string
	CaseList []HttpTestCase
}

type HttpTestCase struct {
	Name          string                               //desc
	EndPoint      string                               //request path, no query included
	Method        string                               //support "P" (POST), "G"(GET)
	BasicAuthFunc func(*HttpTestCase) (string, string) //support nil (no auth) / func
	ContentType   string                               //support "", "J" (application/json; charset=utf-8), "F"(application/x-www-form-urlencoded)
	Body          interface{}                          //support nil, string, utee.J , url.Values(auto encode) , io.Reader
	Status        int                                  //http Status Code
	Res           interface{}                          //support nil (no assert action), string, utee.J, utee.Ast (Ast assert against response json)
}

//Json Assert Object
type Ast struct {
	Path   string      //jsonq path example "data.items.0.pre" , see github.com/jmoiron/jsonq for detail
	Result interface{} //expected result
	Tp     string      //result type, support "I"(int), "S"(string) , "B"(bool) , "E"(regexp match)
}

func (p HttpTestSuite) Exec(handler http.Handler, t *testing.T) {
	for _, op := range p.CaseList {
		op.Name = fmt.Sprint(p.Name, "-", op.Name)
		op.Exec(handler, t)
	}
}

func (p HttpTestCase) GetBody() io.Reader {
	if p.Body == nil || p.Body == "" {
		return nil
	}
	switch p.Body.(type) {
	case string:
		return strings.NewReader(p.Body.(string))
	case J:
		b, _ := json.Marshal(p.Body.(J))
		return bytes.NewReader(b)
	case url.Values:
		return strings.NewReader(p.Body.(url.Values).Encode())
	case io.Reader:
		return p.Body.(io.Reader)
	default:
		panic("unknown body type")
	}
}

func (p HttpTestCase) GetMethod() string {
	switch p.Method {
	case "P":
		return "POST"
	case "G":
		return "GET"
	default:
		log.Fatalln("unknown method", p.Method)
		return ""
	}
}

func (p HttpTestCase) GetEndPoint() string {
	s := p.EndPoint
	if p.GetMethod() == "GET" && p.GetBody() != nil {
		b, err := ioutil.ReadAll(p.GetBody())
		Chk(err)
		s += "?" + string(b)
	}
	return s
}

func (p HttpTestCase) GetContentType() string {
	switch p.ContentType {
	case "J":
		return ContentTypeJson
	case "F":
		return ContentTypeForm
	default:
		return ""
	}
}

func (p HttpTestCase) BuildRequest() *http.Request {
	body := p.GetBody()
	if p.Method == "G" {
		body = nil
	}
	req, _ := http.NewRequest(p.GetMethod(), p.GetEndPoint(), body)
	if p.GetContentType() != "" {
		req.Header.Set("Content-Type", p.GetContentType())
	}
	if p.BasicAuthFunc != nil {
		username, passwd := p.BasicAuthFunc(&p)
		if username != "" {
			req.SetBasicAuth(username, passwd)
		}
	}
	return req
}

func (p HttpTestCase) Exec(handler http.Handler, t *testing.T) {
	req := p.BuildRequest()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, p.Status, w.Code, p.Name)
	if p.Res == nil {
		return
	}
	assertResponse(p.Res, p.Name, t, w)
}

func assertResponse(obj interface{}, desc string, t *testing.T, w *httptest.ResponseRecorder) {
	switch obj.(type) {
	case string:
		s := obj.(string)
		assert.Equal(t, s, w.Body.String(), desc)
	case J:
		j := obj.(J)
		b, _ := json.Marshal(j)
		assert.Equal(t, string(b), w.Body.String(), desc)
	case Ast:
		var j map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &j); err != nil {
			assert.Fail(t, "bad json", "%s, %s", desc, w.Body.String())
			return
		}
		ast := obj.(Ast)
		switch ast.Tp {
		case "I":
			n, err := jsonq.NewQuery(j).Int(strings.Split(ast.Path, ".")...)
			if err != nil {
				assert.Fail(t, "jsonq path err", "%s, %v", ast.Path, err)
				return
			}
			assert.Equal(t, ast.Result, n, desc)
		case "S":
			s, err := jsonq.NewQuery(j).String(strings.Split(ast.Path, ".")...)
			if err != nil {
				assert.Fail(t, "jsonq path err", "%s, %v", ast.Path, err)
				return
			}
			assert.Equal(t, ast.Result, s, desc)
		case "B":
			s, err := jsonq.NewQuery(j).Bool(strings.Split(ast.Path, ".")...)
			if err != nil {
				assert.Fail(t, "jsonq path err", "%s, %v", ast.Path, err)
				return
			}
			assert.Equal(t, ast.Result, s, desc)
		case "E":
			s, err := jsonq.NewQuery(j).String(strings.Split(ast.Path, ".")...)
			if err != nil {
				assert.Fail(t, "jsonq path err", "%s, %v", ast.Path, err)
				return
			}
			m, err := regexp.MatchString(ast.Result.(string), s)
			if err != nil {
				assert.Fail(t, "match fail", "%s, %v", ast.Result, err)
				return
			}
			assert.Equal(t, true, m, s, fmt.Sprintln(desc, ast.Result, s))
		default:
			log.Fatal("unkown ast type", ast.Tp)
		}
	case []interface{}:
		for _, tmp := range obj.([]interface{}) {
			assertResponse(tmp, desc, t, w)
		}
	default:
		log.Fatal("unkown res type", obj)
	}
}
