package service

import (
	"golang.org/x/net/context"
	"ms/card/pkg/common"
	"ms/card/pkg/contract"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/repository"
	"ms/card/pkg/telemetry/jaeger"
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
		Operation             repository.Operations
	}

	Transaction struct {
		TransactionOpts
	}
)

func NewTransaction(opts TransactionOpts) *Transaction {
	return &Transaction{opts}
}

func (t *Transaction) Create(ctx context.Context, request *contract.TransactionRequest) (*entity.Transaction, error) {
	ctx, span := jaeger.Span(ctx)
	defer span.End()

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

	operation, err := t.Operation.FindByID(ctx, request.Operation)
	if err != nil {
		t.Logger.Errorf("t.OperationType.FindByID failed with %s\n", err)
		return nil, err
	}

	if err := t.AccountService.UpdateLimit(ctx, account, request.Amount, operation.Debit); err != nil {
		t.Logger.Errorf("t.AccountService.UpdateLimit failed with %s\n", err)
		return nil, err
	}

	amount := common.Abs(request.Amount)
	if operation.Debit {
		amount = -amount
	}

	transaction, err := t.TransactionRepository.Create(ctx, entity.Transaction{
		Account:   request.Account,
		Type:      request.Operation,
		Amount:    amount,
		CreatedAt: time.Now(),
	})

	if err != nil {
		_ = t.AccountService.UpdateLimit(ctx, account, request.Amount, false)
		t.Logger.Errorf("t.TransactionRepository.Create failed with %s\n", err)
		return nil, err
	}

	return transaction, nil
}
