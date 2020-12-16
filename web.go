package utee

import (
	"encoding/json"
	"io"
	"strings"
)

type J map[string]interface{}

func (p J) ToString() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (p J) ToReader() io.Reader {
	b, _ := json.Marshal(p)
	return strings.NewReader(string(b))
}
