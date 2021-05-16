package service

import (
	e "errors"

	"github.com/google/uuid"
	ent "github.com/lrweck/go-sampleapi/entity"
	rep "github.com/lrweck/go-sampleapi/repository"
	"github.com/pkg/errors"
)

var (
	ErrAccountNotFound    = e.New("account not found")
	ErrAccountInvalid     = e.New("invalid account")
	ErrTransactionInvalid = e.New("invalid transaction")
	ErrInsufficientLimit  = e.New("insufficient account limit")
)

type ApiService struct {
	repo rep.ApiRepository
}

// NewApiService creates a new API service, with which is is possible to call a repository
func NewApiService(repo rep.ApiRepository) *ApiService {
	return &ApiService{repo}
}

// FindAccount fetches a Account struct from the underlying storage, by passing a account ID
func (r *ApiService) FindAccount(acc_id uuid.UUID) (*ent.Account, error) {

	return r.repo.FindAccount(acc_id)
}

// FindAccount persists an account to storage
func (r *ApiService) StoreAccount(acc *ent.Account) error {
	acc.AccountID = uuid.New()
	if acc.AccountLimit == 0 {
		acc.AccountLimit = 200
	}
	return r.repo.StoreAccount(acc)
}

// FindAccount persists an account to storage
func (r *ApiService) StoreTransaction(tx *ent.Transactions) error {

	tx.TransactionID = uuid.New()

	acc, err := r.repo.FindAccount(tx.AccountID)

	if err != nil {
		return errors.Wrap(ErrAccountNotFound, "repository.Transaction.Store")
	}

	normalizeOperationAmount(tx.OpeTypeID, &tx.Amount)

	newLimit := tx.Amount + acc.AccountLimit

	if newLimit < 0 {
		return ErrInsufficientLimit
	}

	tx.NewAccountLimit = newLimit

	return r.repo.StoreTransaction(tx)
}

func normalizeOperationAmount(opeType int, amount *float64) {

	switch opeType {
	case 1, 2, 3:
		if *amount > 0 {
			*amount = -*amount
		}
	case 4:
		if *amount < 0 {
			*amount = -*amount
		}
	}
}
