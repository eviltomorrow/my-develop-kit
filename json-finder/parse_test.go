package finder

import (
	"io/ioutil"
	"testing"
)

func TestGetKey(t *testing.T) {
	buf, err := ioutil.ReadFile("data.json")
	if err != nil {
		t.Fatalf("Error: %v\r\n", err)
	}
	var key = "instanceList.disk.type"
	result, err := GetKey(string(buf), key)
	t.Logf("Error: %v\r\n", err)
	t.Logf("Result: %v\r\n", result)
}
