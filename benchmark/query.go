package benchmark

import (
	"context"
	"log"
	"tdb_benchmarking/db"
	"time"
)

func Process(connStr string, c chan []string, result chan []int64, quit chan int) {
	ctx := context.Background()
	conn, err := db.NewConnection(ctx, connStr)

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(ctx)

	var durations []int64

	for {
		select {
		case query := <-c:
			start := time.Now()

			//nolint:errcheck // ignore the result
			conn.ExecSelect(ctx, query[0], query[1:])

			durations = append(durations, int64(time.Since(start)))
		case <-quit:
			result <- durations
			return
		}
	}
}
