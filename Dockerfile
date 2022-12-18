FROM golang:1.19-alpine AS build
WORKDIR /go/src/timescale-bench
COPY . .
RUN go build -o /go/bin/timescale-bench cmd/main.go

FROM alpine:latest
COPY --from=build /go/bin/timescale-bench /go/bin/timescale-bench
COPY data data
CMD ["/go/bin/timescale-bench"]
