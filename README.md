Timescale Assigment
===================
High level goal: To implement a command line tool that can be used to benchmark SELECT query performance across multiple workers/clients against a TimescaleDB instance. The tool should take as its input a CSV file (whose format is specified below) and a flag to specify the number of concurrent workers. After processing all the queries specified by the parameters in the CSV file, the tool should output a summary with the following stats:

- number of queries processed,
- total processing time across all queries,
- the minimum query time (for a single query),
- the median query time,
- the average query time,
- and the maximum query time.

Build
-----
To build the command run:

  $  go build -o ./timescale-bench cmd/main.go

Usage
-----
The command has the following options:

```
$ ./timescale-bench -h

Usage of ./timescale-bench:
  -data string
      File with query data (default "data/query_params.csv")
  -database string
      Timescale database (default "homework")
  -host string
      Timescale host (default "192.168.1.36")
  -password string
      Timescale password (default "password")
  -port int
      Timescale port (default 5432)
  -username string
      Timescale user (default "postgres")
  -workers int
      Number of workers (default 20)
```

When running the command, it shows:

```
$ ./timescale-bench
Number of queries: 200
Total: 722.753625 ms
Average: 3.613768 ms
Median: 2.223958 ms
Min: 1.672334 ms
Max: 21.533833 ms
```

Run With Docker
---------------
To run the command with Docker, requires docker 20.10.21 and docker-compose v2.13.0

There are the scripts create_db.sh and drop_db.sh to create, import and delete the database in docker compose.

Follow the instructions:

```
# build the docker compose
$ docker-compose -p timescale-bench_compose -f docker-compose.yml build

# start the container
$ docker-compose -p timescale-bench_compose -f docker-compose.yml up -d

# create the database and import data
$ ./create_db.sh

# run the command
$ docker-compose -p timescale-bench_compose -f docker-compose.yml run go /go/bin/timescale-bench -host timescale
```
