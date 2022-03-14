package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	Transaction struct {
		ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
		Account   uint      `json:"account_id" gorm:"type:integer;column:account_id"`
		Type      uint      `json:"operation_type_id" gorm:"type:integer;column:operation_type_id"`
		Amount    int64     `json:"amount" gorm:"type:integer;column:amount"`
		EventDate time.Time `json:"event_date" gorm:"type:timestamp;column:event_date"`
	}

	TransactionCollection struct {
		Total        float64        `json:"total"`
		Transactions []*Transaction `json:"transactions"`
	}
)

func (t *Transaction) TableName() string {
	return "transaction"
}

func (t *TransactionCollection) Sum() {
	values := make([]decimal.Decimal, 0)
	for _, transaction := range t.Transactions {
		values = append(values, decimal.NewFromInt(transaction.Amount))
	}

	sum := decimal.Sum(decimal.NewFromFloat(0), values...)
	total, _ := sum.Float64()

	t.Total = total
}
