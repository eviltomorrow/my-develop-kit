package tool

// Person p
type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Asset Asset  `json:"asset"`
}

// Asset a
type Asset struct {
	Stocks float64 `json:"stocks"`
	Houses int64   `json:"houses"`
}
