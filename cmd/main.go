package main

import (
	"fmt"
	"log"
	"os"

	"tdb_benchmarking/db"
)

func main() {
	log.Println("Up and running")

	conn, ctx := db.Connection()
	defer conn.Close(ctx)

	var greeting string
	err := conn.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
