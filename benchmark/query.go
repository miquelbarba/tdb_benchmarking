package benchmark

import (
	"context"
	"fmt"
	"os"
	"sync"
	"tdb_benchmarking/db"
	"time"
)

func Process(wg *sync.WaitGroup, channel int, c chan []string, quit chan int) {
	defer wg.Done()

	ctx := context.Background()
	conn, err := db.NewConnection("postgres://postgres:password@192.168.1.36:5432/homework", ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(ctx)

	var durations []time.Duration

	for {
		select {
		case query := <-c:
			start := time.Now()

			conn.ExecSelect(
				ctx,
				"SELECT bucket, max, min FROM cpu_usage_summary_minute WHERE host = $1 AND bucket BETWEEN $2 AND $3",
				query[0],
				query[1],
				query[2],
			)

			durations = append(durations, time.Since(start))
		case <-quit:
			fmt.Printf("quit channel %d: len: %d - avg: %d total: %d min: %d, max: %d, median: %d\n", channel,
				len(durations),
				Average(durations),
				Total(durations),
				Min(durations),
				Max(durations),
				Median(durations))
			return
		}
	}
}
