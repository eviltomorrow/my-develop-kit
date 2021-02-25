package sort

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	var data = []int{2, 31, 13, 4, 13, 16, 23, 5, 4, 3, 2, 1}
	QuickSort(data, 0, len(data)-1)
	t.Logf("Data: %v\r\n", data)
}
