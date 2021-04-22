package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

var (
	dbPool *pgx.ConnPool
)

func GetConn() (*pgx.ConnPool, error) {
	if dbPool == nil {

		pool, err := setupDBConn()
		if err != nil {
			log.Fatalf("Connection to DB failed: %v\n", err)
			return nil, err
		}
		dbPool = pool
	}

	return dbPool, nil
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
		AcquireTimeout: 5 * time.Second,
	}

	pool, err := pgx.NewConnPool(pgxConnPool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to DB failed: %v\n", err)
		return nil, err
	}
	return pool, nil
}

func DBCreateAccount(p PostAccount) (uuid.UUID, error) {
	pool, _ := GetConn()

	conn, err := pool.Acquire()

	if err != nil {
		return uuid.Nil, err
	}
	defer pool.Release(conn)

	sql := `INSERT INTO ACCOUNTS (AccountID,DocNumber) VALUES ($1::uuid, $2::text);`

	id := uuid.New()
	if _, err := conn.Exec(sql, id, p.DocNumber); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func DBCreateTransaction(p PostTransaction) error {
	pool, _ := GetConn()

	conn, err := pool.Acquire()

	if err != nil {
		return err
	}

	defer pool.Release(conn)

	_, err = DBGetAccount(p.AccountID)

	if err == pgx.ErrNoRows {
		return err
	}

	// Podemos confiar que o modelo tem as regras de refencia válidas para OperationType e Account.
	// Como não há Login, nem autenticação, qualquer usuário pode inserir transações pra qualquer account
	sql := `INSERT INTO Transactions (TransactionID,AccountID,OperationTypeID,Amount) VALUES
		   ($1::uuid, $2::uuid, $3::smallint, $4::numeric);`

	p.Amount = normalizeOperationAmount(p.OpeType, p.Amount)

	if _, err := conn.Exec(sql, uuid.New(), p.AccountID, p.OpeType, p.Amount); err != nil {
		return err
	}

	return nil
}

func DBGetAccount(aid uuid.UUID) (string, error) {
	pool, _ := GetConn()

	conn, err := pool.Acquire()

	if err != nil {
		return "", err
	}
	defer pool.Release(conn)

	sql := `SELECT DocNumber FROM ACCOUNTS WHERE AccountID = $1::uuid;`

	// fmt.Printf("SELECT DocNumber FROM ACCOUNTS WHERE AccountID = '%v'::uuid;\n", aid)

	docNumber := ""
	err = conn.QueryRow(sql, aid).Scan(&docNumber)
	return docNumber, err
}

// Transações de tipo compra e saque são registradas com valor negativo, enquanto
// transações de pagamento são registradas com valor positivo.
//	1-COMPRA A VISTA
//	2-COMPRA PARCELADA
//	3-SAQUE
//	4-PAGAMENTO
func normalizeOperationAmount(opeType int, amount float64) float64 {

	switch opeType {
	case 1, 2, 3:
		if amount > 0 {
			return -amount
		}
	case 4:
		if amount < 0 {
			return -amount
		}
	}

	return amount
}
