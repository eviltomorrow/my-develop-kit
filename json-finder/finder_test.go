package finder

import "testing"

func TestBuildKeyTree(t *testing.T) {
	var keys = []string{"a.b.c.d", "a.e.f.g", "a.d.f"}
	tree, err := BuildKeyTree(keys)
	if err != nil {
		t.Fatalf("Error: %v\r\n", err)
	}

	t.Logf("Tree: %v", tree)
}
