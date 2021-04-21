package db

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

var (
	dbPool *pgx.ConnPool
)

func GetConn() (*pgx.Conn, error) {
	if dbPool == nil {

		pool, err := setupDBConn()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection to DB failed: %v\n", err)
			panic(1)
		}
		dbPool = pool
	}

	return dbPool.Acquire()
}

func setupDBConn() (*pgx.ConnPool, error) {
	pgxConfig := pgx.ConnConfig{
		Host:     "localhost", // pegar endereço do network do docker
		Database: "sampleapi",
		User:     "sampleapi",
		Password: "supersafe",
	}

	pgxConnPool := pgx.ConnPoolConfig{
		ConnConfig:     pgxConfig,
		MaxConnections: 5,
		AcquireTimeout: 5,
	}

	pool, err := pgx.NewConnPool(pgxConnPool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to DB failed: %v\n", err)
		return nil, err
	}
	return pool, nil
}
