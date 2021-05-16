package repository

import (
	"github.com/google/uuid"
	"github.com/lrweck/go-sampleapi/entity"
)

type ApiRepository interface {
	FindAccount(acc_id uuid.UUID) (*entity.Account, error)
	StoreAccount(acc *entity.Account) error
	StoreTransaction(tx *entity.Transactions) error
}

type repo struct {
	apiRepo ApiRepository
}

func NewApiRepository(rep ApiRepository) ApiRepository {
	return &repo{
		apiRepo: rep,
	}
}

func (r *repo) FindAccount(acc_id uuid.UUID) (*entity.Account, error) {
	return r.apiRepo.FindAccount(acc_id)
}
func (r *repo) StoreAccount(acc *entity.Account) error {
	return r.apiRepo.StoreAccount(acc)
}
func (r *repo) StoreTransaction(tx *entity.Transactions) error {
	return r.apiRepo.StoreTransaction(tx)
}
