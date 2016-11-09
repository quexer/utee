package utee

import (
	"testing"
)

func TestCleanURI(t *testing.T) {
	data := map[string]string{
		"/api/show/101834":                                                "/api/show/:num",
		"/api/attentionCountry/add":                                       "/api/attentionCountry/add",
		"/api/act/57e48711e4d0f314382b45c9/news?since=0&until=0&count=10": "/api/act/:hash/news",
		"/api/img/53f6e4219b8d723f8b00002a/138x138":                       "/api/img/:hash/138x138",
	}

	for k, v := range data {
		if CleanURI(k) != v {
			t.Error("clean fail", k, CleanURI(k))
		}
	}
}
