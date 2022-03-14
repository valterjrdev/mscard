package repository

import (
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"ms/card/pkg/common"

	"ms/card/pkg/persistence"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
)

var (
	ErrAccountCreate              = xerrors.New("failed to create new account")
	ErrAccountCreateAlreadyExists = xerrors.New("account already exists")
	ErrAccountCreateNotFound      = xerrors.New("account not found")
	ErrAccountFindByID            = xerrors.New("failed fetch the account")
)

const (
	UniqueKeyCodeConstraint = "23505"
)

type (
	Accounts interface {
		Create(ctx context.Context, structure entity.Account) (*entity.Account, error)
		FindByID(ctx context.Context, id uint) (*entity.Account, error)
		FindAll(ctx context.Context, filters filter.AccountCollection) ([]*entity.Account, error)
	}

	Account struct {
		logger  common.Logger
		adapter *gorm.DB
	}
)

func NewAccount(logger common.Logger, adapter *gorm.DB) *Account {
	return &Account{
		adapter: adapter,
		logger:  logger,
	}
}

func (a *Account) Create(ctx context.Context, structure entity.Account) (*entity.Account, error) {
	tx := a.adapter.WithContext(ctx)
	if result := tx.Create(&structure); result.Error != nil {
		a.logger.Errorf("tx.Create() failed with %s\n", result.Error)
		var err *pgconn.PgError
		if xerrors.As(result.Error, &err) && err.Code == UniqueKeyCodeConstraint {
			return nil, ErrAccountCreateAlreadyExists
		}

		return nil, ErrAccountCreate
	}

	return &structure, nil
}

func (a *Account) FindByID(ctx context.Context, id uint) (*entity.Account, error) {
	var account entity.Account

	tx := a.adapter.WithContext(ctx)
	if result := tx.Select([]string{"id", "document_number"}).First(&account, id); result.Error != nil {
		a.logger.Errorf("tx.First() failed with %s\n", result.Error)
		if xerrors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAccountCreateNotFound
		}

		return nil, ErrAccountFindByID
	}

	return &account, nil
}

func (a *Account) FindAll(ctx context.Context, filters filter.AccountCollection) ([]*entity.Account, error) {
	accounts := make([]*entity.Account, 0)
	tx := a.adapter.WithContext(ctx)
	find := tx.Scopes(filters.Filter(), persistence.Paginator(filters.Page, filters.Size)).Select([]string{
		"id",
		"document_number",
	}).Find(&accounts)

	return accounts, find.Error
}
