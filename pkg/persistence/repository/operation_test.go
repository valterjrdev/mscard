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

func TestOperationRepository_Create(t *testing.T) {
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
		INSERT INTO "operation" ("description","debit") 
		VALUES ($1,$2) 
		RETURNING "id"
	`)).WithArgs("COMPRA A VISTA", true).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(uint(1)),
	)
	dbmock.ExpectCommit()

	OperationRepository := NewOperation(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operation, err := OperationRepository.Create(ctx, entity.Operation{
		Description: "COMPRA A VISTA",
		Debit:       true,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint(1), Operation.ID)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationRepository_Create_Persist_Error(t *testing.T) {
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
	dbmock.ExpectQuery("^INSERT INTO \"operation\"(.+)$").WillReturnError(ErrOperationCreate)
	dbmock.ExpectRollback()

	OperationRepository := NewOperation(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operation, err := OperationRepository.Create(ctx, entity.Operation{
		Description: "COMPRA A VISTA",
		Debit:       true,
	})

	assert.Nil(t, Operation)
	assert.EqualError(t, err, ErrOperationCreate.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationRepository_Create_Persist_ValidateUniqueKeyConstraint_Error(t *testing.T) {
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
	dbmock.ExpectQuery("^INSERT INTO \"operation\"(.+)$").WillReturnError(&pgconn.PgError{
		Code:    UniqueKeyCodeConstraint,
		Message: ErrOperationCreateAlreadyExists.Error(),
	})

	dbmock.ExpectRollback()

	OperationRepository := NewOperation(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operation, err := OperationRepository.Create(ctx, entity.Operation{
		Description: "COMPRA A VISTA",
		Debit:       true,
	})

	assert.Nil(t, Operation)
	assert.EqualError(t, err, ErrOperationCreateAlreadyExists.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","description","debit"
		FROM "operation" 
		WHERE "operation"."id" = $1 
		ORDER BY "operation"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "description", "debit"}).AddRow(uint(1), "COMPRA A VISTA", true))

	OperationRepository := NewOperation(logger, gormdb)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operation, err := OperationRepository.FindByID(ctx, 1)
	assert.NotNil(t, Operation)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), Operation.ID)
	assert.Equal(t, "COMPRA A VISTA", Operation.Description)
	assert.Equal(t, true, Operation.Debit)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationRepository_FindByID_Error(t *testing.T) {
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
		SELECT "id","description","debit"
		FROM "operation" 
		WHERE "operation"."id" = $1 
		ORDER BY "operation"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnError(ErrOperationFindByID)

	OperationRepository := NewOperation(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operation, err := OperationRepository.FindByID(ctx, 1)
	assert.Nil(t, Operation)
	assert.EqualError(t, err, ErrOperationFindByID.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationRepository_FindByID_RecordNotFound_Error(t *testing.T) {
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
		SELECT "id","description","debit"
		FROM "operation" 
		WHERE "operation"."id" = $1 
		ORDER BY "operation"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnError(gorm.ErrRecordNotFound)

	OperationRepository := NewOperation(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operation, err := OperationRepository.FindByID(ctx, 1)
	assert.Nil(t, Operation)
	assert.EqualError(t, err, ErrOperationCreateNotFound.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOperationRepository_Collection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","description","debit"
		FROM "operation"
		LIMIT 10
	`)).WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).
		AddRow(uint(1), "COMPRA A VISTA").
		AddRow(uint(2), "COMPRA PARCELADA"),
	)

	OperationRepository := NewOperation(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	Operations, err := OperationRepository.FindAll(ctx, filter.OperationCollection{})
	assert.Len(t, Operations, 2)
	assert.NoError(t, err)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
