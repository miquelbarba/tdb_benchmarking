package benchmark

import (
	"testing"

	"tdb_benchmarking/benchmark"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	connStr := "postgres://postgres:password@localhost:5432"

	c := make(chan []string)
	quit := make(chan int)
	result := make(chan []int64)

	go benchmark.Process(connStr, c, result, quit)

	c <- []string{"SELECT 1 WHERE $1 = $2 AND $2 = $3", "a", "a", "a"}
	c <- []string{"SELECT 1 WHERE $1 = $2 AND $2 = $3", "b", "b", "b"}
	quit <- 1

	res := <-result

	assert.Equal(t, 2, len(res))
	assert.Greater(t, res[0], int64(0))
	assert.Greater(t, res[1], int64(0))
}
