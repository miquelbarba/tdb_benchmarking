package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connection() (*pgx.Conn, context.Context) {
	ctx := context.Background()
	connStr := "postgres://postgres:password@192.168.1.36:5432/homework"

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, ctx
}
