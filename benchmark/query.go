package benchmark

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const ConnectionError = 4

func Process(connStr string, c chan []string, result chan []int64, quit chan int) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		fmt.Println(err)
		os.Exit(ConnectionError)
	}

	defer conn.Close(ctx)

	var durations []int64

	for {
		select {
		case query := <-c:
			start := time.Now()

			rows, err := conn.Query(ctx, query[0], query[1], query[2], query[3])

			if err != nil {
				fmt.Printf("Error executing %s: %s\n", query[1:], err)
			}

			durations = append(durations, int64(time.Since(start)))

			rows.Close()
		case <-quit:
			result <- durations
			return
		}
	}
}
