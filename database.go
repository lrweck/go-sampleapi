package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

// Var global pra manter o estado
var (
	dbPool    *pgx.ConnPool
	ErrNoRows = pgx.ErrNoRows
)

// Retorna uma conexão do pool.
// Inicia a conexão ao BD, caso ainda não tenha sido feita
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

// Inicializa os parametros referentes ao pool de conexões
func setupDBConn() (*pgx.ConnPool, error) {
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	pgxConfig := pgx.ConnConfig{
		Host:     os.Getenv("PG_HOST"),
		Port:     uint16(port),
		Database: os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
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

// Func pra criar accounts.
// É chamada pelo handler HandlerCreateAccount
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

// Func pra criar transactions.
// É chamada pelo handler HandlerCreateTransaction
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

// Func pra retornar uma account.
// É chamada pelo handler HandlerGetAccount e pelo
// HandlerCreateTransaction (para validar a existencia da conta)
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
