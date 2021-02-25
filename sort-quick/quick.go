package sort

// QuickSort quick sort
func QuickSort(data []int, begin, end int) []int {
	if begin > end {
		return data
	}

	leftIndex, rightIndex := begin, end
	var pivotIndex = leftIndex
	var pivot = data[pivotIndex]

	// 挖坑
	for leftIndex < rightIndex {
		if pivotIndex == leftIndex {
			if data[rightIndex] < pivot {
				data[pivotIndex] = data[rightIndex]
				pivotIndex = rightIndex
				leftIndex++
			} else {
				rightIndex--
			}
		}

		if pivotIndex == rightIndex {
			if data[leftIndex] > pivot {
				data[pivotIndex] = data[leftIndex]
				pivotIndex = leftIndex
				rightIndex--
			} else {
				leftIndex++
			}
		}

	}
	data[pivotIndex] = pivot
	data = QuickSort(data, begin, pivotIndex-1)
	data = QuickSort(data, pivotIndex+1, end)
	return data
}
