package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
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

func parseCommandLine() (numWorkers int, file, connStr string) {
	workers := flag.Int("workers", DefaultWorkers, fmt.Sprintf("Number of workers (default %d)", DefaultWorkers))
	fileName := flag.String("data", DefaultFile, fmt.Sprintf("File with query data (default %s)", DefaultFile))

	host := flag.String("host", DefaultHost, fmt.Sprintf("Timescale host (default %s)", DefaultHost))
	port := flag.Int("port", DefaultPort, fmt.Sprintf("Timescale port (default %d)", DefaultPort))
	user := flag.String("username", DefaultUser, fmt.Sprintf("Timescale user (default %s)", DefaultUser))
	password := flag.String("password", DefaultPassword, fmt.Sprintf("Timescale password (default %s)", DefaultPassword))
	database := flag.String("database", DefaultDatabase, fmt.Sprintf("Timescale database (default %s)", DefaultDatabase))

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

	data := readCSV(fileName)
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
