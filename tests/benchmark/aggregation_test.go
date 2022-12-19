package benchmark

import (
	"testing"

	"tdb_benchmarking/benchmark"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTotal(t *testing.T) {
	assert.Equal(t, int64(0), benchmark.Total([]int64{}))
	assert.Equal(t, int64(33), benchmark.Total([]int64{3, 4, 5, 21}))
}

func TestAverage(t *testing.T) {
	assert.Equal(t, decimal.NewFromInt(0), benchmark.Average([]int64{}))

	expected, _ := decimal.NewFromString("8.25")
	assert.True(t, expected.Equal(benchmark.Average([]int64{3, 4, 5, 21})))
}

func TestMin(t *testing.T) {
	assert.Equal(t, int64(0), benchmark.Min([]int64{}))
	assert.Equal(t, int64(3), benchmark.Min([]int64{3, 4, 5, 21}))
	assert.Equal(t, int64(0), benchmark.Min([]int64{54, 0, 21}))
}

func TestMax(t *testing.T) {
	assert.Equal(t, int64(0), benchmark.Max([]int64{}))
	assert.Equal(t, int64(21), benchmark.Max([]int64{3, 4, 5, 21}))
	assert.Equal(t, int64(54), benchmark.Max([]int64{54, 0, 21}))
}

func TestMedian(t *testing.T) {
	assert.Equal(t, decimal.NewFromInt(0), benchmark.Median([]int64{}))

	expected, _ := decimal.NewFromString("4.5")
	assert.True(t, expected.Equal(benchmark.Median([]int64{3, 4, 5, 21})))

	assert.Equal(t, decimal.NewFromInt(3), benchmark.Median([]int64{3, 1, 8}))
}
