package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ms/card/pkg/common"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
	"regexp"
	"testing"
	"time"
)

func TestOperationTypeRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectBegin()
	dbmock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "operation_type" ("description","negative") 
		VALUES ($1,$2) 
		RETURNING "id"
	`)).WithArgs("COMPRA A VISTA", true).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(uint(1)),
	)
	dbmock.ExpectCommit()

	OperationTypeRepository := NewOperationType(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationType, err := OperationTypeRepository.Create(ctx, entity.OperationType{
		Description: "COMPRA A VISTA",
		Negative:    true,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint(1), operationType.ID)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationTypeRepository_Create_Persist_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectBegin()
	dbmock.ExpectQuery("^INSERT INTO \"operation_type\"(.+)$").WillReturnError(ErrOperationTypeCreate)
	dbmock.ExpectRollback()

	OperationTypeRepository := NewOperationType(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationType, err := OperationTypeRepository.Create(ctx, entity.OperationType{
		Description: "COMPRA A VISTA",
		Negative:    true,
	})

	assert.Nil(t, operationType)
	assert.EqualError(t, err, ErrOperationTypeCreate.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationTypeRepository_Create_Persist_ValidateUniqueKeyConstraint_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectBegin()
	dbmock.ExpectQuery("^INSERT INTO \"operation_type\"(.+)$").WillReturnError(&pgconn.PgError{
		Code:    UniqueKeyCodeConstraint,
		Message: ErrOperationTypeCreateAlreadyExists.Error(),
	})

	dbmock.ExpectRollback()

	OperationTypeRepository := NewOperationType(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationType, err := OperationTypeRepository.Create(ctx, entity.OperationType{
		Description: "COMPRA A VISTA",
		Negative:    true,
	})

	assert.Nil(t, operationType)
	assert.EqualError(t, err, ErrOperationTypeCreateAlreadyExists.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationTypeRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","description","negative"
		FROM "operation_type" 
		WHERE "operation_type"."id" = $1 
		ORDER BY "operation_type"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "description", "negative"}).AddRow(uint(1), "COMPRA A VISTA", true))

	OperationTypeRepository := NewOperationType(logger, gormdb)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationType, err := OperationTypeRepository.FindByID(ctx, 1)
	assert.NotNil(t, operationType)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), operationType.ID)
	assert.Equal(t, "COMPRA A VISTA", operationType.Description)
	assert.Equal(t, true, operationType.Negative)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationTypeRepository_FindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","description","negative"
		FROM "operation_type" 
		WHERE "operation_type"."id" = $1 
		ORDER BY "operation_type"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnError(ErrOperationTypeFindByID)

	OperationTypeRepository := NewOperationType(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationType, err := OperationTypeRepository.FindByID(ctx, 1)
	assert.Nil(t, operationType)
	assert.EqualError(t, err, ErrOperationTypeFindByID.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationTypeRepository_FindByID_RecordNotFound_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","description","negative"
		FROM "operation_type" 
		WHERE "operation_type"."id" = $1 
		ORDER BY "operation_type"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnError(gorm.ErrRecordNotFound)

	OperationTypeRepository := NewOperationType(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationType, err := OperationTypeRepository.FindByID(ctx, 1)
	assert.Nil(t, operationType)
	assert.EqualError(t, err, ErrOperationTypeCreateNotFound.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationTypeRepository_Collection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","description","negative"
		FROM "operation_type"
		LIMIT 10
	`)).WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).
		AddRow(uint(1), "COMPRA A VISTA").
		AddRow(uint(2), "COMPRA PARCELADA"),
	)

	OperationTypeRepository := NewOperationType(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	operationTypes, err := OperationTypeRepository.FindAll(ctx, filter.OperationTypeCollection{})
	assert.Len(t, operationTypes, 2)
	assert.NoError(t, err)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
