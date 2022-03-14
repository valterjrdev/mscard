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
	ErrOperationTypeCreate              = xerrors.New("failed to create new operation type")
	ErrOperationTypeCreateAlreadyExists = xerrors.New("operation type already exists")
	ErrOperationTypeCreateNotFound      = xerrors.New("operation type not found")
	ErrOperationTypeFindByID            = xerrors.New("failed fetch operation type")
)

type (
	OperationTypes interface {
		Create(ctx context.Context, structure entity.OperationType) (*entity.OperationType, error)
		FindByID(ctx context.Context, id uint) (*entity.OperationType, error)
		FindAll(ctx context.Context, filters filter.OperationTypeCollection) ([]*entity.OperationType, error)
	}

	OperationType struct {
		logger  common.Logger
		adapter *gorm.DB
	}
)

func NewOperationType(logger common.Logger, adapter *gorm.DB) *OperationType {
	return &OperationType{
		adapter: adapter,
		logger:  logger,
	}
}

func (a *OperationType) Create(ctx context.Context, structure entity.OperationType) (*entity.OperationType, error) {
	tx := a.adapter.WithContext(ctx)
	if result := tx.Create(&structure); result.Error != nil {
		a.logger.Errorf("tx.Create() failed with %s\n", result.Error)
		var err *pgconn.PgError
		if xerrors.As(result.Error, &err) && err.Code == UniqueKeyCodeConstraint {
			return nil, ErrOperationTypeCreateAlreadyExists
		}

		return nil, ErrOperationTypeCreate
	}

	return &structure, nil
}

func (a *OperationType) FindByID(ctx context.Context, id uint) (*entity.OperationType, error) {
	var operationType entity.OperationType

	tx := a.adapter.WithContext(ctx)
	if result := tx.Select([]string{"id", "description", "negative"}).First(&operationType, id); result.Error != nil {
		a.logger.Errorf("tx.First() failed with %s\n", result.Error)
		if xerrors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOperationTypeCreateNotFound
		}

		return nil, ErrOperationTypeFindByID
	}

	return &operationType, nil
}

func (a *OperationType) FindAll(ctx context.Context, filters filter.OperationTypeCollection) ([]*entity.OperationType, error) {
	operationTypes := make([]*entity.OperationType, 0)
	tx := a.adapter.WithContext(ctx)
	find := tx.Scopes(filters.Filter(), persistence.Paginator(filters.Page, filters.Size)).Select([]string{
		"id",
		"description",
		"negative",
	}).Find(&operationTypes)

	return operationTypes, find.Error
}
