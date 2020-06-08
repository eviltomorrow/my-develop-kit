package tool

import "testing"

func TestFindObejctByTag(t *testing.T) {
	var p = Person{
		Name: "shepard",
		Age:  20,
		Asset: Asset{
			Stocks: 20000,
			Houses: 3,
		},
	}

	rt, rv, err := FindObjectByTag([]string{"asset", "stocks"}, &p)
	t.Logf("error: %v", err)
	t.Logf("rt: %v", rt)
	t.Logf("rv: %v, Canset: %v", rv, rv.CanSet())
}
