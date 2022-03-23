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
	ErrOperationCreate              = xerrors.New("failed to create new operation")
	ErrOperationCreateAlreadyExists = xerrors.New("operation already exists")
	ErrOperationCreateNotFound      = xerrors.New("operation not found")
	ErrOperationFindByID            = xerrors.New("failed fetch operation")
)

type (
	Operations interface {
		Create(ctx context.Context, structure entity.Operation) (*entity.Operation, error)
		FindByID(ctx context.Context, id uint) (*entity.Operation, error)
		FindAll(ctx context.Context, filters filter.OperationCollection) ([]*entity.Operation, error)
	}

	Operation struct {
		logger  common.Logger
		adapter *gorm.DB
	}
)

func NewOperation(logger common.Logger, adapter *gorm.DB) *Operation {
	return &Operation{
		adapter: adapter,
		logger:  logger,
	}
}

func (a *Operation) Create(ctx context.Context, structure entity.Operation) (*entity.Operation, error) {
	tx := a.adapter.WithContext(ctx)
	if result := tx.Create(&structure); result.Error != nil {
		a.logger.Errorf("tx.Create() failed with %s\n", result.Error)
		var err *pgconn.PgError
		if xerrors.As(result.Error, &err) && err.Code == UniqueKeyCodeConstraint {
			return nil, ErrOperationCreateAlreadyExists
		}

		return nil, ErrOperationCreate
	}

	return &structure, nil
}

func (a *Operation) FindByID(ctx context.Context, id uint) (*entity.Operation, error) {
	var operation entity.Operation

	tx := a.adapter.WithContext(ctx)
	if result := tx.Select([]string{"id", "description", "debit"}).First(&operation, id); result.Error != nil {
		a.logger.Errorf("tx.First() failed with %s\n", result.Error)
		if xerrors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOperationCreateNotFound
		}

		return nil, ErrOperationFindByID
	}

	return &operation, nil
}

func (a *Operation) FindAll(ctx context.Context, filters filter.OperationCollection) ([]*entity.Operation, error) {
	operations := make([]*entity.Operation, 0)
	tx := a.adapter.WithContext(ctx)
	find := tx.Scopes(filters.Filter(), persistence.Paginator(filters.Page, filters.Size)).Select([]string{
		"id",
		"description",
		"debit",
	}).Find(&operations)

	return operations, find.Error
}
