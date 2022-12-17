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
const ConnStr = "postgres://postgres:password@192.168.1.36:5432/homework"
const NumWorkers = 20
const FileName = "data/query_params.csv"

func buildWorkers(numWorkers int) (channels []chan []string, result chan []int64, quit chan int) {
	quit = make(chan int)
	result = make(chan []int64)

	for i := 0; i < numWorkers; i++ {
		channels = append(channels, make(chan []string))
		go benchmark.Process(ConnStr, channels[i], result, quit)
	}

	return channels, result, quit
}

func quitWorkers(quit chan int, result chan []int64, numWorkers int) []int64 {
	var totalDuration []int64

	for i := 0; i < numWorkers; i++ {
		quit <- 0

		durations := <-result
		totalDuration = append(totalDuration, durations...)
	}

	return totalDuration
}

func readCSV(fileName string) [][]string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Panicln(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		log.Panicln(err)
	}

	return data
}

func main() {
	data := readCSV(FileName)
	channels, result, quit := buildWorkers(NumWorkers)

	for i := 1; i < len(data); i++ {
		numChannel, _ := strconv.Atoi(strings.Split(data[i][0], "_")[1])
		channels[numChannel] <- append([]string{Query}, data[i]...)
	}

	totalDuration := quitWorkers(quit, result, NumWorkers)

	fmt.Printf(
		"Number of queries: %d\n"+
			"Average time: %d\n"+
			"Total time: %d\n"+
			"Min: %d\n"+
			"Max: %d\n"+
			"Median: %d\n",
		len(totalDuration),
		benchmark.Average(totalDuration),
		benchmark.Total(totalDuration),
		benchmark.Min(totalDuration),
		benchmark.Max(totalDuration),
		benchmark.Median(totalDuration))
}
