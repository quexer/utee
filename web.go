package utee

import "encoding/json"

type J map[string]interface{}

func (p J) ToString() string {
	b, _ := json.Marshal(p)
	return string(b)
}
