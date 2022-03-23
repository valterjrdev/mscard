package repository

import (
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"ms/card/pkg/common"
	"ms/card/pkg/persistence"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
)

var (
	ErrTransactionCreate = xerrors.New("failed to create new transaction")
)

type (
	Transactions interface {
		Create(ctx context.Context, structure entity.Transaction) (*entity.Transaction, error)
		FindAll(ctx context.Context, filters filter.TransactionCollection) (*entity.TransactionCollection, error)
	}

	Transaction struct {
		logger  common.Logger
		adapter *gorm.DB
	}
)

func NewTransaction(logger common.Logger, adapter *gorm.DB) *Transaction {
	return &Transaction{
		adapter: adapter,
		logger:  logger,
	}
}

func (a *Transaction) Create(ctx context.Context, structure entity.Transaction) (*entity.Transaction, error) {
	tx := a.adapter.WithContext(ctx)
	if result := tx.Create(&structure); result.Error != nil {
		a.logger.Errorf("tx.Create() failed with %s\n", result.Error)
		return nil, ErrTransactionCreate
	}

	return &structure, nil
}

func (a *Transaction) FindAll(ctx context.Context, filters filter.TransactionCollection) (*entity.TransactionCollection, error) {
	transactions := make([]*entity.Transaction, 0)
	tx := a.adapter.WithContext(ctx)
	find := tx.Scopes(filters.Filter(), persistence.Paginator(filters.Page, filters.Size)).Select([]string{
		"id",
		"account_id",
		"operation_id",
		"amount",
		"event_date",
	}).Find(&transactions)

	collection := &entity.TransactionCollection{Transactions: transactions}
	collection.Sum()
	return collection, find.Error
}
