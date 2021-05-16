package service

import (
	"github.com/lrweck/go-sampleapi/entity"
)

type AccountSerializer interface {
	Decode(input []byte) (*entity.Account, error)
	Encode(input *entity.Account) ([]byte, error)
}

type TransactionSerializer interface {
	Decode(input []byte) (*entity.Transactions, error)
	Encode(input *entity.Transactions) ([]byte, error)
}
