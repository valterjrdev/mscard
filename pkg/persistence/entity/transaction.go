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
		EventDate time.Time `json:"event_date" gorm:"type:timestamp;column:event_date"`
	}

	TransactionCollection struct {
		Total        int64          `json:"total"`
		Transactions []*Transaction `json:"transactions"`
	}
)

func (t *Transaction) TableName() string {
	return TransactionTableName
}

func (t *TransactionCollection) Sum() {
	values := int64(0)
	for _, transaction := range t.Transactions {
		values += transaction.Amount
	}

	t.Total = values
}
