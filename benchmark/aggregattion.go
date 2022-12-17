package benchmark

import (
	"sort"
	"time"
)

func Total(arr []time.Duration) int64 {
	total := int64(0)
	for _, duration := range arr {
		total += int64(duration)
	}

	return total
}

func Average(arr []time.Duration) int64 {
	if len(arr) == 0 {
		return 0
	}

	return Total(arr) / int64(len(arr))
}

func Min(arr []time.Duration) int64 {
	if len(arr) == 0 {
		return 0
	}

	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}

	return int64(min)
}

func Max(arr []time.Duration) int64 {
	if len(arr) == 0 {
		return 0
	}

	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return int64(max)
}

func Median(arr []time.Duration) int64 {
	dataCopy := make([]time.Duration, len(arr))
	copy(dataCopy, arr)

	sort.Slice(dataCopy, func(i, j int) bool { return dataCopy[i] < dataCopy[j] })

	l := len(dataCopy)

	if l == 0 {
		return 0
	}

	if l%2 == 0 {
		return int64((dataCopy[l/2-1] + dataCopy[l/2]) / 2) //nolint:gomnd // no magic number
	}

	return int64(dataCopy[l/2])
}
