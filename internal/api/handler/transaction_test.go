package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
	"ms/card/pkg/persistence/repository"
	"ms/card/pkg/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlerTransaction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	mockTransactionService := service.NewMockTransactions(ctrl)
	mockTransactionService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Transaction{
		ID:        1,
		Account:   1,
		Type:      4,
		Amount:    10020,
		EventDate: time.Now(),
	}, nil)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"account_id":1,"operation_id":4,"amount":10020}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewTransaction(TransactionOpts{
		TransactionService: mockTransactionService,
	})

	if assert.NoError(t, h.Create(echo.New().NewContext(req, rec))) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"id":1,"account_id":1,"operation_id":4,"amount":10020,"event_date":"2022-03-12T01:02:03.000000004Z"}`, rec.Body.String())
	}
}

func TestHandlerTransaction_Create_Persist_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	mockTransactionService := service.NewMockTransactions(ctrl)
	mockTransactionService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, repository.ErrTransactionCreate)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"account_id":1,"operation_id":4,"amount":10020}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewTransaction(TransactionOpts{
		TransactionService: mockTransactionService,
	})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=failed to create new transaction")
}

func TestHandlerTransaction_Create_BindRequest_Erro(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"a`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewTransaction(TransactionOpts{})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=code=400, message=unexpected EOF, internal=unexpected EOF")
}

func TestHandlerTransaction_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	collection := &entity.TransactionCollection{
		Data: []*entity.Transaction{
			{
				ID:        1,
				Account:   1,
				Type:      4,
				Amount:    -10000,
				EventDate: time.Now(),
			},
			{
				ID:        2,
				Account:   1,
				Type:      4,
				Amount:    5000,
				EventDate: time.Now(),
			},
		},
	}

	collection.Sum()

	mockTranscationRepository := repository.NewMockTransactions(ctrl)
	mockTranscationRepository.EXPECT().FindAll(gomock.Any(), filter.TransactionCollection{}).Return(collection, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewTransaction(TransactionOpts{
		TransactionRepository: mockTranscationRepository,
	})

	if assert.NoError(t, h.FindAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `
		{
			"balance": -5000,
			"data": [
				{"id":1,"account_id":1,"operation_id":4,"amount":-10000,"event_date":"2022-03-12T01:02:03.000000004Z"},
				{"id":2,"account_id":1,"operation_id":4,"amount":5000,"event_date":"2022-03-12T01:02:03.000000004Z"}
			]
		}
		`, rec.Body.String())
	}
}

func TestHandlerTransaction_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time { return time.Date(2022, time.March, 12, 1, 2, 3, 4, time.UTC) })
	assert.NoError(t, err)
	defer func() {
		if err := patch.Unpatch(); err != nil {
			t.Error(err)
		}
	}()

	mockTranscationRepository := repository.NewMockTransactions(ctrl)
	mockTranscationRepository.EXPECT().FindAll(gomock.Any(), filter.TransactionCollection{}).Return(nil, errors.New("err find all"))

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewTransaction(TransactionOpts{
		TransactionRepository: mockTranscationRepository,
	})

	assert.EqualError(t, h.FindAll(c), "code=400, message=err find all")
}
