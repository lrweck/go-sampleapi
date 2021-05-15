package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	AccountID    uuid.UUID `json:"Account_ID"`
	DocNumber    string    `json:"Document_Number"`
	AccountLimit float64   `json:"AvailableCreditLimit"`
}

type OperationTypes struct {
	OpeTypeID   int    `json:"OperationType_ID"`
	Description string `json:"Description0"`
}

type Transactions struct {
	TransactionID   uuid.UUID `json:"Transaction_ID"`
	AccountID       uuid.UUID `json:"Account_ID"`
	OpeTypeID       int       `json:"OperationType_ID"`
	Amount          float64   `json:"Amount"`
	EventDate       time.Time `json:"EventDate"`
	NewAccountLimit float64
}
