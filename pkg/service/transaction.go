package service

import (
	"golang.org/x/net/context"
	"ms/card/pkg/common"
	"ms/card/pkg/contract"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/repository"
	"time"
)

type (
	Transactions interface {
		Create(ctx context.Context, request *contract.TransactionRequest) (*entity.Transaction, error)
	}

	TransactionOpts struct {
		Logger                common.Logger
		AccountService        Accounts
		TransactionRepository repository.Transactions
		AccountRepository     repository.Accounts
		OperationType         repository.OperationTypes
	}

	Transaction struct {
		TransactionOpts
	}
)

func NewTransaction(opts TransactionOpts) *Transaction {
	return &Transaction{opts}
}

func (t *Transaction) Create(ctx context.Context, request *contract.TransactionRequest) (*entity.Transaction, error) {
	if err := request.Validate(); err != nil {
		t.Logger.Errorf("request.Validate() failed with %s\n", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	account, err := t.AccountRepository.FindByID(ctx, request.Account)
	if err != nil {
		t.Logger.Errorf("t.AccountRepository.FindByID failed with %s\n", err)
		return nil, err
	}

	operationType, err := t.OperationType.FindByID(ctx, request.Type)
	if err != nil {
		t.Logger.Errorf("t.OperationType.FindByID failed with %s\n", err)
		return nil, err
	}

	if err := t.AccountService.UpdateLimit(ctx, account, request.Amount, operationType.Negative); err != nil {
		t.Logger.Errorf("t.AccountService.UpdateLimit failed with %s\n", err)
		return nil, err
	}

	amount := common.Abs(request.Amount)
	if operationType.Negative {
		amount = -amount
	}

	transaction, err := t.TransactionRepository.Create(ctx, entity.Transaction{
		Account:   request.Account,
		Type:      request.Type,
		Amount:    amount,
		EventDate: time.Now(),
	})

	if err != nil {
		if err := t.AccountService.UpdateLimit(ctx, account, request.Amount, false); err != nil {
			t.Logger.Warnf("t.AccountService.UpdateLimit failed with %s\n", err)
		}

		t.Logger.Errorf("t.TransactionRepository.Create failed with %s\n", err)
		return nil, err
	}

	return transaction, nil
}
