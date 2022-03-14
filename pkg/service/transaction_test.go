package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"ms/card/pkg/common"
	"ms/card/pkg/contract"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/repository"
	"testing"
)

func TestServiceTransaction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountEntity := &entity.Account{
		ID:       1,
		Document: "56077053074",
		Limit:    2000,
	}
	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(mockAccountEntity, nil)

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.OperationType{
		ID:          1,
		Description: "COMPRA A VISTA",
		Negative:    true,
	}, nil)

	mockTransactionRepository := repository.NewMockTransactions(ctrl)
	mockTransactionRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Transaction{
		ID:      1,
		Account: 1,
		Type:    1,
		Amount:  1000,
	}, nil)

	accountServiceMock := NewMockAccounts(ctrl)
	accountServiceMock.EXPECT().UpdateLimit(gomock.Any(), mockAccountEntity, int64(1000), true).Return(nil)

	transactionService := NewTransaction(TransactionOpts{
		AccountService:        accountServiceMock,
		OperationType:         mockOperationTypeRepository,
		AccountRepository:     mockAccountRepository,
		TransactionRepository: mockTransactionRepository,
	})

	transaction, err := transactionService.Create(context.Background(), &contract.TransactionRequest{
		Account: 1,
		Type:    1,
		Amount:  1000,
	})
	assert.NotNil(t, transaction)
	assert.Nil(t, err)
}

func TestServiceTransaction_Create_Limit_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := common.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockAccountEntity := &entity.Account{
		ID:       1,
		Document: "56077053074",
		Limit:    100,
	}
	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(mockAccountEntity, nil)

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.OperationType{
		ID:          1,
		Description: "COMPRA A VISTA",
		Negative:    true,
	}, nil)

	accountServiceMock := NewMockAccounts(ctrl)
	accountServiceMock.EXPECT().UpdateLimit(gomock.Any(), mockAccountEntity, int64(1000), true).Return(ErrLimitExceeded)

	transactionService := NewTransaction(TransactionOpts{
		Logger:            mockLogger,
		AccountService:    accountServiceMock,
		OperationType:     mockOperationTypeRepository,
		AccountRepository: mockAccountRepository,
	})

	transaction, err := transactionService.Create(context.Background(), &contract.TransactionRequest{
		Account: 1,
		Type:    1,
		Amount:  1000,
	})
	assert.Nil(t, transaction)
	assert.EqualError(t, err, ErrLimitExceeded.Error())
}

func TestServiceTransaction_Create_Validate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := common.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	transactionService := NewTransaction(TransactionOpts{
		Logger: mockLogger,
	})

	transaction, err := transactionService.Create(context.Background(), &contract.TransactionRequest{})
	assert.Nil(t, transaction)
	assert.EqualError(t, err, "account_id: cannot be blank; amount: cannot be blank; operation_type_id: cannot be blank.")
}

func TestServiceTransaction_Create_Persist_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := common.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockAccountEntity := &entity.Account{
		ID:       1,
		Document: "56077053074",
		Limit:    2000,
	}
	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(mockAccountEntity, nil)

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.OperationType{
		ID:          1,
		Description: "COMPRA A VISTA",
		Negative:    true,
	}, nil)

	accountServiceMock := NewMockAccounts(ctrl)
	accountServiceMock.EXPECT().UpdateLimit(gomock.Any(), mockAccountEntity, int64(1000), true).Return(nil)
	accountServiceMock.EXPECT().UpdateLimit(gomock.Any(), mockAccountEntity, int64(1000), false).Return(nil)

	mockTransactionRepository := repository.NewMockTransactions(ctrl)
	mockTransactionRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, repository.ErrTransactionCreate)

	transactionService := NewTransaction(TransactionOpts{
		Logger:                mockLogger,
		AccountService:        accountServiceMock,
		OperationType:         mockOperationTypeRepository,
		AccountRepository:     mockAccountRepository,
		TransactionRepository: mockTransactionRepository,
	})

	transaction, err := transactionService.Create(context.Background(), &contract.TransactionRequest{
		Account: 1,
		Type:    1,
		Amount:  1000,
	})
	assert.Nil(t, transaction)
	assert.EqualError(t, err, repository.ErrTransactionCreate.Error())
}

func TestServiceTransaction_Create_Persist_Account_NotFound_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := common.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(nil, repository.ErrAccountCreateNotFound)

	transactionService := NewTransaction(TransactionOpts{
		Logger:            mockLogger,
		AccountRepository: mockAccountRepository,
	})

	transaction, err := transactionService.Create(context.Background(), &contract.TransactionRequest{
		Account: 1,
		Type:    1,
		Amount:  1000,
	})
	assert.Nil(t, transaction)
	assert.EqualError(t, err, repository.ErrAccountCreateNotFound.Error())
}

func TestServiceTransaction_Create_Persist_OperationType_NotFound_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := common.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.Account{
		ID:       1,
		Document: "56077053074",
	}, nil)

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(nil, repository.ErrOperationTypeCreateNotFound)

	transactionService := NewTransaction(TransactionOpts{
		Logger:            mockLogger,
		OperationType:     mockOperationTypeRepository,
		AccountRepository: mockAccountRepository,
	})

	transaction, err := transactionService.Create(context.Background(), &contract.TransactionRequest{
		Account: 1,
		Type:    1,
		Amount:  1000,
	})
	assert.Nil(t, transaction)
	assert.EqualError(t, err, repository.ErrOperationTypeCreateNotFound.Error())
}
