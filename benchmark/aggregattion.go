package benchmark

import "sort"

func Total(arr []int64) int64 {
	total := int64(0)
	for _, duration := range arr {
		total += duration
	}

	return total
}

func Average(arr []int64) int64 {
	if len(arr) == 0 {
		return 0
	}

	return Total(arr) / int64(len(arr))
}

func Min(arr []int64) int64 {
	if len(arr) == 0 {
		return 0
	}

	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}

	return min
}

func Max(arr []int64) int64 {
	if len(arr) == 0 {
		return 0
	}

	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}

func Median(arr []int64) int64 {
	dataCopy := make([]int64, len(arr))
	copy(dataCopy, arr)

	sort.Slice(dataCopy, func(i, j int) bool { return dataCopy[i] < dataCopy[j] })

	l := len(dataCopy)

	if l == 0 {
		return 0
	}

	if l%2 == 0 {
		//nolint:gomnd // no magic number
		return (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	}

	return dataCopy[l/2]
}

func ToMilliseconds(num int64) float64 {
	//nolint:gomnd // no magic number
	return float64(num) / 1e6
}
