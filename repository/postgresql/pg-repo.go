package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	pool "github.com/jackc/pgx/v4/pgxpool"
	ent "github.com/lrweck/go-sampleapi/entity"
	rep "github.com/lrweck/go-sampleapi/repository"
	serv "github.com/lrweck/go-sampleapi/service"
	"github.com/pkg/errors"
)

const (
	maxPoolConns    = 5
	minPoolConns    = 1
	maxConnIdleTime = time.Minute
)

type pgRepo struct {
	conn    *pool.Pool
	timeout time.Duration
}

func newPgClient(pgURL string, timeout int) (*pool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	config, err := pool.ParseConfig(pgURL)

	config.MaxConns = maxPoolConns
	config.MinConns = minPoolConns
	config.MaxConnIdleTime = maxConnIdleTime

	if err != nil {
		return nil, err
	}

	connPool, err := pool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = connPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return connPool, nil
}

// NewPGRepo reates a new PostgreSQl repository to store and consume data.
func NewPGRepo(pgURL string, timeout int) (rep.ApiRepository, error) {
	repo := &pgRepo{
		timeout: time.Duration(timeout) * time.Second,
	}
	conn, err := newPgClient(pgURL, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPGRepo")
	}
	repo.conn = conn
	return repo, nil
}

// FindAccount locates the account and returns an Account struct
func (r *pgRepo) FindAccount(acc_id uuid.UUID) (*ent.Account, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	acc := &ent.Account{
		AccountID: acc_id,
	}
	sql := `SELECT DocNumber,AvailableCreditLimit FROM ACCOUNTS WHERE AccountID = $1::uuid LIMIT 1;`

	err := r.conn.QueryRow(ctx, sql, acc_id).Scan(&acc.DocNumber, &acc.AccountLimit)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(serv.ErrAccountNotFound, "repository.Account.Find")
		}
		return nil, errors.Wrap(err, "repository.Account.Find")
	}
	return acc, nil
}

// StoreAccount persists the account struct to the database
func (r *pgRepo) StoreAccount(acc *ent.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	sql := `INSERT INTO ACCOUNTS (AccountID,DocNumber,availablecreditlimit) VALUES ($1::uuid, $2::text, $3::numeric);`
	_, err := r.conn.Exec(ctx, sql, acc.AccountID, acc.DocNumber, acc.AccountLimit)

	if err != nil {
		return errors.Wrap(err, "repository.Account.Store")
	}
	return nil
}

// StoreTransaction persists the transaction struct to the database.
// Updates the account limit reflecting the changes
func (r *pgRepo) StoreTransaction(transac *ent.Transactions) error {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	sql := `INSERT INTO Transactions (TransactionID,AccountID,OperationTypeID,Amount) VALUES
	($1::uuid, $2::uuid, $3::smallint, $4::numeric);`

	tx, err := r.conn.Begin(ctx)

	if err != nil {
		fmt.Printf("1-%s\n", err)
		return errors.Wrap(err, "repository.Transaction.Store")
	}

	_, err = tx.Exec(ctx, sql, transac.TransactionID, transac.AccountID, transac.OpeTypeID, transac.Amount)

	if err != nil {
		fmt.Printf("2-%s\ntx: %+v", err, transac)
		tx.Rollback(ctx)
		return errors.Wrap(err, "repository.Transaction.Store")
	}

	if err = updateAccountLimit(ctx, tx, transac.AccountID, transac.NewAccountLimit); err != nil {
		fmt.Printf("3-%s\n", err)
		tx.Rollback(ctx)
		return errors.Wrap(err, "repository.Transaction.Store")
	}

	tx.Commit(ctx)

	return nil
}

func updateAccountLimit(ctx context.Context, tx pgx.Tx, acc uuid.UUID, amount float64) error {

	sql := `UPDATE ACCOUNTS SET AvailableCreditLimit = $2 WHERE AccountID = $1::uuid;`

	if _, err := tx.Exec(ctx, sql, acc, amount); err != nil {
		return err
	}

	return nil
}
