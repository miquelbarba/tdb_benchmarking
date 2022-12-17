package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"tdb_benchmarking/benchmark"
)

const Query = "SELECT bucket, max, min FROM cpu_usage_summary_minute WHERE host = $1 AND bucket BETWEEN $2 AND $3"
const connStr = "postgres://postgres:password@192.168.1.36:5432/homework"
const NumWorkers = 20

func main() {
	quit := make(chan int)
	result := make(chan []int64)

	var channels [NumWorkers]chan []string
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan []string)
		go benchmark.Process(connStr, channels[i], result, quit)
	}

	f, err := os.Open("data/query_params.csv")
	if err != nil {
		log.Panicln(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	//nolint:errcheck // we don't use header
	csvReader.Read()

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panicln(err)
		}

		numChannel, err := strconv.Atoi(strings.Split(rec[0], "_")[1])

		if err != nil {
			log.Panicln(err)
		}

		channels[numChannel] <- append([]string{Query}, rec...)
	}

	for i := 0; i < len(channels); i++ {
		quit <- 0
	}

	var totalDuration []int64

	for i := 0; i < len(channels); i++ {
		durations := <-result
		totalDuration = append(totalDuration, durations...)
	}

	fmt.Printf("Number of queries: %d\nAverage: %d\nTotal time: %d\nMin: %d\nMax: %d\nMedian: %d\n",
		len(totalDuration),
		benchmark.Average(totalDuration),
		benchmark.Total(totalDuration),
		benchmark.Min(totalDuration),
		benchmark.Max(totalDuration),
		benchmark.Median(totalDuration))
}
