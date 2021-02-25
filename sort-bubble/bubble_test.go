package sort

import "testing"

func TestBubbleSort(t *testing.T) {
	var data = []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	data = BubbleSort(data)
	t.Logf("data: %v", data)
}
