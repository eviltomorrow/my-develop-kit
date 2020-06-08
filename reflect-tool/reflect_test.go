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

	err = SetString("name", "captain", &p)
	t.Logf("SetString: %v, value: %s", err, p.String())

	err = SetInt64("asset.houses", 100000000, &p)
	t.Logf("SetInt64: %v, value: %s", err, p.String())

	rt, rv, err = FindObjectByTag([]string{"child"}, &p)
	t.Logf("error: %v", err)
	t.Logf("rt: %v", rt)
	t.Logf("rv: %v, Canset: %v", rv, rv.CanSet())
	SetSliceString("child", []string{"1", "b"}, &p)
	t.Logf("SetSliceString: %v, value: %s", err, p.String())

}
