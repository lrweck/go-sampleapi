package json

import (
	"encoding/json"

	"github.com/lrweck/go-sampleapi/entity"
	"github.com/pkg/errors"
)

// Redirect struct to add methods to.
type Account struct{}
type Transactions struct{}

// Decode bytes into an Account struct via json
func (r *Account) Decode(input []byte) (*entity.Account, error) {
	acc := &entity.Account{}
	if err := json.Unmarshal(input, acc); err != nil {
		return nil, errors.Wrap(err, "serializer.json.Account.Decode")
	}
	return acc, nil
}

// Encode an Account struct to bytes
func (r *Account) Encode(input *entity.Account) ([]byte, error) {
	return genericEncode(input)
}

// Encode a Transactions struct to bytes
func (r *Transactions) Encode(input *entity.Transactions) ([]byte, error) {
	return genericEncode(input)
}

// Decode bytes into a Transactions struct via json
func (r *Transactions) Decode(input []byte) (*entity.Transactions, error) {
	tx := &entity.Transactions{}
	if err := json.Unmarshal(input, tx); err != nil {
		return nil, errors.Wrap(err, "serializer.json.Transactions.Decode")
	}
	return tx, nil
}

func genericEncode(obj interface{}) ([]byte, error) {
	rawMsg, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.json.Encode")
	}
	return rawMsg, nil
}
