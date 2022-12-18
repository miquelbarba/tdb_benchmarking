package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"tdb_benchmarking/benchmark"
)

// -workers 40

const Query = "SELECT bucket, max, min FROM cpu_usage_summary_minute WHERE host = $1 AND bucket BETWEEN $2 AND $3"
const ConnStr = "postgres://postgres:password@192.168.1.36:5432/homework"
const DefaultNumWorkers = 20
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

func assignWorker(identifier string, numWorkers int, mappingWorkers map[string]int, currentWorker *int) int {
	worker, ok := mappingWorkers[identifier]
	if ok {
		return worker
	}

	mappingWorkers[identifier] = *currentWorker % numWorkers
	(*currentWorker)++

	fmt.Printf("%s -> %d \n", identifier, mappingWorkers[identifier])

	return mappingWorkers[identifier]
}

func parseCommandLine() int {
	numWorkers := flag.Int("workers", DefaultNumWorkers, "Number of workers")
	flag.Parse()

	return *numWorkers
}

func main() {
	numWorkers := parseCommandLine()

	currentWorker := 0
	mappingWorkers := make(map[string]int)

	data := readCSV(FileName)
	channels, result, quit := buildWorkers(numWorkers)

	for i := 1; i < len(data); i++ {
		numWorker := assignWorker(data[i][0], numWorkers, mappingWorkers, &currentWorker)
		channels[numWorker] <- append([]string{Query}, data[i]...)
	}

	totalDuration := quitWorkers(quit, result, numWorkers)

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
