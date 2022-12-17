package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"tdb_benchmarking/benchmark"
)

const Query = "SELECT bucket, max, min FROM cpu_usage_summary_minute WHERE host = $1 AND bucket BETWEEN $2 AND $3"
const connStr = "postgres://postgres:password@192.168.1.36:5432/homework"
const NumWorkers = 20
const File = "data/query_params.csv"

func main() {
	quit := make(chan int)
	result := make(chan []int64)

	var channels [NumWorkers]chan []string
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan []string)
		go benchmark.Process(connStr, channels[i], result, quit)
	}

	f, err := os.Open(File)
	if err != nil {
		log.Panicln(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		log.Panicln(err)
	}

	for i := 1; i < len(data); i++ {
		numChannel, _ := strconv.Atoi(strings.Split(data[i][0], "_")[1])
		channels[numChannel] <- append([]string{Query}, data[i]...)
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
