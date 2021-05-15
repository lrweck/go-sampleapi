package main

import (
	"errors"
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
	dbPool               *pgx.ConnPool
	ErrNoRows            = pgx.ErrNoRows
	ErrInsufficientLimit = errors.New("transação inválida. saldo da conta insuficiente")
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

	sql := `INSERT INTO ACCOUNTS (AccountID,DocNumber,availablecreditlimit) VALUES ($1::uuid, $2::text, $3::numeric);`

	id := uuid.New()

	if p.AccountLimit <= 0 {
		p.AccountLimit = 200
	}

	if _, err := conn.Exec(sql, id, p.DocNumber, p.AccountLimit); err != nil {
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

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	acc, err := DBGetAccount(p.AccountID)

	if err == pgx.ErrNoRows {
		return err
	}

	p.Amount = normalizeOperationAmount(p.OpeType, p.Amount)

	newLimit := p.Amount + acc.AccountLimit
	if newLimit < 0 {
		return ErrInsufficientLimit
	}

	// Podemos confiar que o modelo tem as regras de refencia válidas para OperationType e Account.
	// Como não há Login, nem autenticação, qualquer usuário pode inserir transações pra qualquer account
	sql := `INSERT INTO Transactions (TransactionID,AccountID,OperationTypeID,Amount) VALUES
		   ($1::uuid, $2::uuid, $3::smallint, $4::numeric);`

	if _, err := tx.Exec(sql, uuid.New(), p.AccountID, p.OpeType, p.Amount); err != nil {
		tx.Rollback()
		return err
	}

	if err = updateAccountLimit(tx, p.AccountID, newLimit); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Erro ao realizar commit: %v\n", err)
		return err
	}

	return nil

}

// Func pra retornar uma account.
// É chamada pelo handler HandlerGetAccount e pelo
// HandlerCreateTransaction (para validar a existencia da conta)
func DBGetAccount(aid uuid.UUID) (*Account, error) {
	pool, _ := GetConn()

	conn, err := pool.Acquire()

	if err != nil {
		return nil, err
	}
	defer pool.Release(conn)

	sql := `SELECT AccountID,DocNumber,AvailableCreditLimit FROM ACCOUNTS WHERE AccountID = $1::uuid;`

	Acc := Account{}

	err = conn.QueryRow(sql, aid).Scan(&Acc.AccountID, &Acc.DocNumber, &Acc.AccountLimit)

	log.Printf("Account: %+v\n", Acc)
	return &Acc, err
}

func updateAccountLimit(tx *pgx.Tx, acc uuid.UUID, amount float64) error {

	sql := `UPDATE ACCOUNTS SET AvailableCreditLimit = $2 WHERE AccountID = $1::uuid;`

	if _, err := tx.Exec(sql, acc, amount); err != nil {
		return err
	}

	return nil
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
