package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"tdb_benchmarking/benchmark"
)

const Query = "SELECT bucket, max, min FROM cpu_usage_summary_minute WHERE host = $1 AND bucket BETWEEN $2 AND $3"

func main() {
	var wg sync.WaitGroup

	wg.Add(20) //nolint:gomnd // no magic number

	quit := make(chan int)

	var channels [20]chan []string
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan []string)
		go benchmark.Process(&wg, i, channels[i], quit)
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

	wg.Wait()
}
