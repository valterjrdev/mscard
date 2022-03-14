package repository

import (
	"errors"
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

func TestAccountRepository_Create(t *testing.T) {
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
		INSERT INTO "account" ("document_number","limit") 
		VALUES ($1,$2) 
		RETURNING "id"
	`)).WithArgs("64715245019", int64(2000)).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow("1"),
	)
	dbmock.ExpectCommit()

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := accountRepository.Create(ctx, entity.Account{
		Document: "64715245019",
		Limit:    2000,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint(1), account.ID)
	assert.Equal(t, int64(2000), account.Limit)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_Create_Persist_Error(t *testing.T) {
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
	dbmock.ExpectQuery("^INSERT INTO \"account\"(.+)$").WillReturnError(ErrAccountCreate)
	dbmock.ExpectRollback()

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := accountRepository.Create(ctx, entity.Account{
		Document: "64715245019",
	})

	assert.Nil(t, account)
	assert.EqualError(t, err, ErrAccountCreate.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_Create_Persist_ValidateUniqueKeyConstraint_Error(t *testing.T) {
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
	dbmock.ExpectQuery("^INSERT INTO \"account\"(.+)$").WillReturnError(&pgconn.PgError{
		Code:    UniqueKeyCodeConstraint,
		Message: ErrAccountCreateAlreadyExists.Error(),
	})

	dbmock.ExpectRollback()

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := accountRepository.Create(ctx, entity.Account{
		Document: "64715245019",
	})

	assert.Nil(t, account)
	assert.EqualError(t, err, ErrAccountCreateAlreadyExists.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","document_number","limit"
		FROM "account" 
		WHERE "account"."id" = $1 
		ORDER BY "account"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "document_number", "limit"}).AddRow(uint(1), "64715245019", int64(2000)))

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := accountRepository.FindByID(ctx, 1)
	assert.NotNil(t, account)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), account.ID)
	assert.Equal(t, "64715245019", account.Document)
	assert.Equal(t, int64(2000), account.Limit)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_FindByID_Error(t *testing.T) {
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
		SELECT "id","document_number","limit"
		FROM "account" 
		WHERE "account"."id" = $1 
		ORDER BY "account"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnError(ErrAccountFindByID)

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := accountRepository.FindByID(ctx, 1)
	assert.Nil(t, account)
	assert.EqualError(t, err, ErrAccountFindByID.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_FindByID_RecordNotFound_Error(t *testing.T) {
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
		SELECT "id","document_number","limit"
		FROM "account" 
		WHERE "account"."id" = $1 
		ORDER BY "account"."id" 
		LIMIT 1
	`)).WithArgs(uint(1)).WillReturnError(gorm.ErrRecordNotFound)

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := accountRepository.FindByID(ctx, 1)
	assert.Nil(t, account)
	assert.EqualError(t, err, ErrAccountCreateNotFound.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_Collection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","document_number","limit"
		FROM "account"
		LIMIT 10
	`)).WillReturnRows(sqlmock.NewRows([]string{"id", "document_number", "limit"}).
		AddRow(uint(1), "64715245019", int64(2000)).
		AddRow(uint(2), "11115245019", int64(2000)),
	)

	accountRepository := NewAccount(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	accounts, err := accountRepository.FindAll(ctx, filter.AccountCollection{})
	assert.Len(t, accounts, 2)
	assert.NoError(t, err)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccount_UpdateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectBegin()
	dbmock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "account" SET "document_number"=$1,"limit"=$2 WHERE "id" = $3`,
	)).WillReturnResult(sqlmock.NewResult(1, 1))
	dbmock.ExpectCommit()

	accountRepository := NewAccount(logger, gormdb)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	accountEntity := &entity.Account{
		ID:       1,
		Document: "64715245019",
		Limit:    2000,
	}

	err = accountRepository.UpdateLimit(ctx, accountEntity)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), accountEntity.ID)
	assert.Equal(t, int64(2000), accountEntity.Limit)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestAccount_UpdateLimit_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	errExpected := errors.New("update err")
	dbmock.ExpectBegin()
	dbmock.ExpectExec(regexp.QuoteMeta(`UPDATE "account" SET "document_number"=$1,"limit"=$2 WHERE "id" = $3`)).WillReturnError(errExpected)
	dbmock.ExpectRollback()

	accountRepository := NewAccount(logger, gormdb)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	accountEntity := &entity.Account{
		ID:       1,
		Document: "64715245019",
		Limit:    2000,
	}

	err = accountRepository.UpdateLimit(ctx, accountEntity)
	assert.EqualError(t, err, errExpected.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
