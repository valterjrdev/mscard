package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_TableName(t *testing.T) {
	transaction := Transaction{}
	assert.Equal(t, TransactionTableName, transaction.TableName())
}

func TestTransactionCollection_TableName(t *testing.T) {
	collection := TransactionCollection{
		Data: []*Transaction{
			{
				Amount: 1000,
			},
			{
				Amount: 2000,
			},
			{
				Amount: 3000,
			},
		},
	}
	collection.Sum()
	assert.Equal(t, int64(6000), collection.Balance)
}
