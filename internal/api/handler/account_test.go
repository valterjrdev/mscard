package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
	"ms/card/pkg/persistence/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerAccount_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Account{
		ID:       1,
		Document: "56077053074",
	}, nil)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"document_number":"56077053074"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewAccount(AccountOpts{
		AccountRepository: mockAccountRepository,
	})

	if assert.NoError(t, h.Create(echo.New().NewContext(req, rec))) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"id":1,"document_number":"56077053074"}`, rec.Body.String())
	}
}

func TestHandlerAccount_Create_Persist_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, repository.ErrAccountCreate)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"document_number":"56077053074"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewAccount(AccountOpts{
		AccountRepository: mockAccountRepository,
	})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=failed to create new account")
}

func TestHandlerAccount_Create_BindRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewAccount(AccountOpts{})
	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=code=400, message=unexpected EOF, internal=unexpected EOF")
}

func TestHandlerAccount_Create_Validate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"document_number":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewAccount(AccountOpts{})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=document_number: cannot be blank.")
}

func TestHandlerAccount_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.Account{
		ID:       1,
		Document: "56077053074",
	}, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetPath(AccountFindByIDPath)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewAccount(AccountOpts{
		AccountRepository: mockAccountRepository,
	})

	if assert.NoError(t, h.FindByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"id":1,"document_number":"56077053074"}`, rec.Body.String())
	}
}

func TestHandlerAccount_FindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(nil, repository.ErrAccountFindByID)

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetPath(AccountFindByIDPath)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewAccount(AccountOpts{
		AccountRepository: mockAccountRepository,
	})

	assert.EqualError(t, h.FindByID(c), "code=400, message=failed fetch the account")
}

func TestHandlerAccount_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindAll(gomock.Any(), filter.AccountCollection{}).Return([]*entity.Account{
		{
			ID:       1,
			Document: "56077053074",
		},
		{
			ID:       2,
			Document: "87756158008",
		},
	}, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewAccount(AccountOpts{
		AccountRepository: mockAccountRepository,
	})

	if assert.NoError(t, h.FindAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `
		[
			{"id":1,"document_number":"56077053074"},
			{"id":2,"document_number":"87756158008"}
		]
		`, rec.Body.String())
	}
}

func TestHandlerAccount_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepository := repository.NewMockAccounts(ctrl)
	mockAccountRepository.EXPECT().FindAll(gomock.Any(), filter.AccountCollection{}).Return(nil, errors.New("err find all"))

	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewAccount(AccountOpts{
		AccountRepository: mockAccountRepository,
	})

	assert.EqualError(t, h.FindAll(c), "code=400, message=err find all")
}
