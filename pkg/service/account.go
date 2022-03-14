package service

import (
	"errors"
	"golang.org/x/net/context"
	"ms/card/pkg/common"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/repository"
)

var (
	ErrLimitExceeded = errors.New("account limit exceeded, operation not allowed")
)

type (
	Accounts interface {
		UpdateLimit(ctx context.Context, account *entity.Account, amount int64, negative bool) error
	}

	AccountOpts struct {
		Logger            common.Logger
		AccountRepository repository.Accounts
	}

	Account struct {
		AccountOpts
	}
)

func NewAccount(opts AccountOpts) *Account {
	return &Account{opts}
}

func (a *Account) UpdateLimit(ctx context.Context, account *entity.Account, amount int64, negative bool) error {
	amount = common.Abs(amount)

	if negative {
		if amount > account.Limit {
			return ErrLimitExceeded
		}

		account.Limit -= amount
		return a.AccountRepository.UpdateLimit(ctx, account)
	}

	account.Limit += amount
	return a.AccountRepository.UpdateLimit(ctx, account)
}
