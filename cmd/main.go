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

func main() {
	var wg sync.WaitGroup
	wg.Add(20)

	quit := make(chan int)

	var chans [20]chan []string
	for i := 0; i < len(chans); i++ {
		chans[i] = make(chan []string)
		go benchmark.Process(&wg, i, chans[i], quit)
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
	csvReader.Read()

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		numChannel, err := strconv.Atoi(strings.Split(rec[0], "_")[1])
		chans[numChannel] <- rec
	}

	for i := 0; i < len(chans); i++ {
		quit <- 0
	}

	wg.Wait()
}
