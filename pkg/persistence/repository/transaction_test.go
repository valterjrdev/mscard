package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ms/card/pkg/common"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
	"regexp"
	"testing"
	"time"
)

func TestTransactionRepository_Create(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

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
		INSERT INTO "transaction" ("account_id","operation_id","amount","event_date") 
		VALUES ($1,$2,$3,$4) 
		RETURNING "id"
	`)).WithArgs(1, 4, 12345, time.Now()).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow("1"),
	)
	dbmock.ExpectCommit()

	transactionRepository := NewTransaction(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	transaction, err := transactionRepository.Create(ctx, entity.Transaction{
		Account:   1,
		Type:      4,
		Amount:    12345,
		EventDate: time.Now(),
	})
	assert.NoError(t, err)
	assert.Equal(t, uint(1), transaction.ID)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionRepository_Create_Persist_Error(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

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
	dbmock.ExpectQuery("^INSERT INTO \"transaction\"(.+)$").WillReturnError(ErrTransactionCreate)
	dbmock.ExpectRollback()

	transactionRepository := NewTransaction(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	account, err := transactionRepository.Create(ctx, entity.Transaction{
		Account:   1,
		Type:      4,
		Amount:    12345,
		EventDate: time.Now(),
	})

	assert.Nil(t, account)
	assert.EqualError(t, err, ErrTransactionCreate.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionRepository_Collection(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	dbmock.ExpectQuery(regexp.QuoteMeta(`SELECT "id","account_id","operation_id","amount","event_date" FROM "transaction" LIMIT 10`)).WillReturnRows(
		sqlmock.NewRows([]string{"id", "account_id", "operation_id", "amount", "event_date"}).
			AddRow(uint(1), uint(1), uint(1), 12340, time.Now()).
			AddRow(uint(2), uint(2), uint(2), 10040, time.Now()).
			AddRow(uint(3), uint(3), uint(3), 10040, time.Now()),
	)

	transactionRepository := NewTransaction(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	collection, err := transactionRepository.FindAll(ctx, filter.TransactionCollection{})
	assert.NoError(t, err)
	assert.Len(t, collection.Transactions, 3)

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionRepository_Collection_Fetch_Error(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := common.NewMockLogger(ctrl)

	mockdb, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()

	gormdb, err := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{})
	assert.NoError(t, err)

	expected := xerrors.New("error fetch collection")
	dbmock.ExpectQuery(regexp.QuoteMeta(`
		SELECT "id","account_id","operation_id","amount","event_date" 
		FROM "transaction" LIMIT 10
	`)).WillReturnError(expected)

	transactionRepository := NewTransaction(logger, gormdb)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	collection, err := transactionRepository.FindAll(ctx, filter.TransactionCollection{})
	assert.Len(t, collection.Transactions, 0)
	assert.EqualError(t, err, expected.Error())

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
