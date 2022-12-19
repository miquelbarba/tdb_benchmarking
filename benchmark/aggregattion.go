package benchmark

import (
	"sort"

	"github.com/shopspring/decimal"
)

func Total(arr []int64) int64 {
	total := int64(0)
	for _, duration := range arr {
		total += duration
	}

	return total
}

func Average(arr []int64) decimal.Decimal {
	if len(arr) == 0 {
		return decimal.NewFromInt(0)
	}

	return decimal.NewFromInt(Total(arr)).Div(decimal.NewFromInt(int64(len(arr))))
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

func Median(arr []int64) decimal.Decimal {
	if len(arr) == 0 {
		return decimal.NewFromInt(0)
	}

	dataCopy := make([]int64, len(arr))
	copy(dataCopy, arr)

	sort.Slice(dataCopy, func(i, j int) bool { return dataCopy[i] < dataCopy[j] })

	l := len(dataCopy)

	midElem := decimal.NewFromInt(dataCopy[l/2])

	if l%2 == 0 {
		//nolint:gomnd // no magic number
		return decimal.NewFromInt(dataCopy[l/2-1]).Add(midElem).Div(decimal.NewFromInt(2))
	}

	return midElem
}

func ToMilliseconds(num int64) float64 {
	//nolint:gomnd // no magic number
	return float64(num) / 1e6
}

func DecimalToMilliseconds(num decimal.Decimal) decimal.Decimal {
	//nolint:gomnd // no magic number
	return num.Div(decimal.NewFromInt(1e6))
}
