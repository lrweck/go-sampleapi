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

type apiService struct {
	repo rep.ApiRepository
}

func NewApiService(repo rep.ApiRepository) rep.ApiRepository {
	return &apiService{repo}
}

func (r *apiService) FindAccount(acc_id uuid.UUID) (*ent.Account, error) {
	return r.repo.FindAccount(acc_id)
}

func (r *apiService) StoreAccount(acc *ent.Account) error {
	acc.AccountID = uuid.New()
	return r.repo.StoreAccount(acc)
}

func (r *apiService) StoreTransaction(tx *ent.Transactions) error {

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

	return amount
}
