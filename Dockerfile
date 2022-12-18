FROM golang:1.19-alpine

RUN apk add --no-cache git

WORKDIR /go/src/timescale-bench

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /go/bin/timescale-bench cmd/main.go

CMD ["/go/bin/timescale-bench"]