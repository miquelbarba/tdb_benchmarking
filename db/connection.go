package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Connection struct {
	DB *pgx.Conn
}

// connStr := "postgres://postgres:password@192.168.1.36:5432/homework"
func NewConnection(connStr string, ctx context.Context) (Connection, error) {
	conn, err := pgx.Connect(ctx, connStr)

	return Connection{DB: conn}, err
}

// SQL query that returns the max cpu usage and min cpu usage of the given hostname
// for every minute in the time range specified by the start time and end time.
// SELECT MAX(usage), MIN(usage) FROM cpu_usage WHERE host = $1 AND ts BETWEEN($2, $3)
func (c *Connection) ExecSelect(ctx context.Context, query string, args ...any) {
	c.DB.Query(ctx, query, args)
}

func (c *Connection) Close(ctx context.Context) {
	c.DB.Close(ctx)
}
