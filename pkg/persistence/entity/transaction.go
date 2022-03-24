package entity

import (
	"time"
)

const (
	TransactionTableName = "transaction"
)

type (
	Transaction struct {
		ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
		Account   uint      `json:"account_id" gorm:"type:integer;column:account_id"`
		Type      uint      `json:"operation_id" gorm:"type:integer;column:operation_id"`
		Amount    int64     `json:"amount" gorm:"type:integer;column:amount"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp without time zone;column:created_at"`
	}

	TransactionCollection struct {
		Balance int64          `json:"balance"`
		Data    []*Transaction `json:"data"`
	}
)

func (t *Transaction) TableName() string {
	return TransactionTableName
}

func (t *TransactionCollection) Sum() {
	values := int64(0)
	for _, transaction := range t.Data {
		values += transaction.Amount
	}

	t.Balance = values
}
