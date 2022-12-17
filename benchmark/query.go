package benchmark

import (
	"context"
	"fmt"
	"log"
	"sync"
	"tdb_benchmarking/db"
	"time"
)

func Process(wg *sync.WaitGroup, channel int, c chan []string, quit chan int) {
	ctx := context.Background()
	conn, err := db.NewConnection(ctx, "postgres://postgres:password@192.168.1.36:5432/homework")

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(ctx)
	defer wg.Done()

	var durations []int64

	for {
		select {
		case query := <-c:
			start := time.Now()

			//nolint:errcheck // ignore the result
			conn.ExecSelect(ctx, query[0], query[1:])

			durations = append(durations, int64(time.Since(start)))
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
