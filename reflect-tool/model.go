package tool

import "encoding/json"

// Person p
type Person struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	Asset     Asset    `json:"asset"`
	ChildName []string `json:"child"`
}

// Asset a
type Asset struct {
	Stocks float64 `json:"stocks"`
	Houses int64   `json:"houses"`
}

func (p *Person) String() string {
	buf, _ := json.Marshal(p)
	return string(buf)
}
