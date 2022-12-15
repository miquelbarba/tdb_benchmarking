package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"tdb_benchmarking/db"
)

func process(wg *sync.WaitGroup, channel int, c chan []string, quit chan int) {
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
				average(durations),
				total(durations),
				min(durations),
				max(durations),
				median(durations))
			return
		}
	}
}

func total(arr []time.Duration) int64 {
	total := int64(0)
	for _, duration := range arr {
		total = total + int64(duration)
	}

	return total
}

func average(arr []time.Duration) int64 {
	if len(arr) == 0 {
		return 0
	}

	return total(arr) / int64(len(arr))
}

func min(arr []time.Duration) int64 {
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

func max(arr []time.Duration) int64 {
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

func median(arr []time.Duration) int64 {
	dataCopy := make([]time.Duration, len(arr))
	copy(dataCopy, arr)

	sort.Slice(dataCopy, func(i, j int) bool { return dataCopy[i] < dataCopy[j] })

	var median int64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = int64((dataCopy[l/2-1] + dataCopy[l/2]) / 2)
	} else {
		median = int64(dataCopy[l/2])
	}

	return median
}

func main() {
	var wg sync.WaitGroup
	wg.Add(20)

	quit := make(chan int)

	var chans [20]chan []string
	for i := 0; i < len(chans); i++ {
		chans[i] = make(chan []string)
		go process(&wg, i, chans[i], quit)
	}

	// open file
	f, err := os.Open("data/query_params.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)

	// header
	csvReader.Read()

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%+v\n", rec)

		numChannel, err := strconv.Atoi(strings.Split(rec[0], "_")[1])
		chans[numChannel] <- rec
	}

	for i := 0; i < len(chans); i++ {
		quit <- 0
	}

	wg.Wait()
}
