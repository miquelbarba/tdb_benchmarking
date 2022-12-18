Timescale Assigment
===================
Cloud Engineer Assignment - Benchmarking

Usage
-----


Run With Docker
---------------
Requires docker 20.10.21 and docker-compose v2.13.0

To run the command with Docker follow the instructions:

  $ docker-compose -p timescale-bench_compose -f docker-compose.yml build
  $ docker-compose -p timescale-bench_compose -f docker-compose.yml up -d
  $ ./create_db.sh

  $ docker-compose -p timescale-bench_compose -f docker-compose.yml run go /go/bin/timescale-bench -host timescale