package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"tdb_benchmarking/benchmark"
)

const Query = "SELECT bucket, max, min FROM cpu_usage_summary_minute WHERE host = $1 AND bucket BETWEEN $2 AND $3"

const DefaultHost = "192.168.1.36"
const DefaultUser = "postgres"
const DefaultPassword = "password"
const DefaultPort = 5432
const DefaultDatabase = "homework"
const DefaultWorkers = 20
const DefaultFile = "data/query_params.csv"

func buildWorkers(numWorkers int, connStr string) (channels []chan []string, result chan []int64, quit chan int) {
	quit = make(chan int)
	result = make(chan []int64)

	for i := 0; i < numWorkers; i++ {
		channels = append(channels, make(chan []string))
		go benchmark.Process(connStr, channels[i], result, quit)
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

func readCSV(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func assignWorker(identifier string, numWorkers int, mappingWorkers map[string]int, currentWorker *int) int {
	worker, ok := mappingWorkers[identifier]
	if ok {
		return worker
	}

	mappingWorkers[identifier] = *currentWorker % numWorkers
	(*currentWorker)++

	return mappingWorkers[identifier]
}

func parseCommandLine() (numWorkers int, file, connStr string) {
	workers := flag.Int("workers", DefaultWorkers, "Number of workers")
	fileName := flag.String("data", DefaultFile, "File with query data")

	host := flag.String("host", DefaultHost, "Timescale host")
	port := flag.Int("port", DefaultPort, "Timescale port")
	user := flag.String("username", DefaultUser, "Timescale user")
	password := flag.String("password", DefaultPassword, "Timescale password")
	database := flag.String("database", DefaultDatabase, "Timescale database")

	flag.Parse()

	numWorkers = *workers
	file = *fileName
	connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", *user, *password, *host, *port, *database)

	return
}

func main() {
	numWorkers, fileName, connStr := parseCommandLine()

	currentWorker := 0
	mappingWorkers := make(map[string]int)

	data, err := readCSV(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	channels, result, quit := buildWorkers(numWorkers, connStr)

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
