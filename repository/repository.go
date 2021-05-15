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
