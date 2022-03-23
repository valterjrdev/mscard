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

func TestHandlerOperation_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationRepository := repository.NewMockOperations(ctrl)
	mockOperationRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Operation{
		ID:          1,
		Description: "COMPRA A VISTA",
		Debit:       true,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"description":"COMPRA A VISTA","debit":"true"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperation(OperationOpts{
		OperationRepository: mockOperationRepository,
	})

	if assert.NoError(t, h.Create(echo.New().NewContext(req, rec))) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"id":1,"description":"COMPRA A VISTA","debit":true}`, rec.Body.String())
	}
}

func TestHandlerOperation_Create_Persist_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationRepository := repository.NewMockOperations(ctrl)
	mockOperationRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, repository.ErrOperationCreate)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"description":"COMPRA A VISTA","debit":"true"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperation(OperationOpts{
		OperationRepository: mockOperationRepository,
	})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=failed to create new operation")
}

func TestHandlerOperation_Create_BindRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperation(OperationOpts{})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=code=400, message=unexpected EOF, internal=unexpected EOF")
}

func TestHandlerOperation_Create_Validate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"description":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperation(OperationOpts{})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=debit: cannot be blank; description: cannot be blank.")
}

func TestHandlerOperation_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationRepository := repository.NewMockOperations(ctrl)
	mockOperationRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.Operation{
		ID:          1,
		Description: "COMPRA A VISTA",
		Debit:       true,
	}, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetPath(AccountFindByIDPath)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewOperation(OperationOpts{
		OperationRepository: mockOperationRepository,
	})

	if assert.NoError(t, h.FindByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"id":1,"description":"COMPRA A VISTA","debit":true}`, rec.Body.String())
	}
}

func TestHandlerOperation_FindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationRepository := repository.NewMockOperations(ctrl)
	mockOperationRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(nil, repository.ErrOperationFindByID)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetPath(AccountFindByIDPath)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewOperation(OperationOpts{
		OperationRepository: mockOperationRepository,
	})

	assert.EqualError(t, h.FindByID(c), "code=400, message=failed fetch operation")
}

func TestHandlerOperation_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationRepository := repository.NewMockOperations(ctrl)
	mockOperationRepository.EXPECT().FindAll(gomock.Any(), filter.OperationCollection{}).Return([]*entity.Operation{
		{
			ID:          1,
			Description: "COMPRA A VISTA",
			Debit:       true,
		},
		{
			ID:          2,
			Description: "PAGAMENTO",
			Debit:       false,
		},
	}, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewOperation(OperationOpts{
		OperationRepository: mockOperationRepository,
	})

	if assert.NoError(t, h.FindAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `
		[
			{"id":1,"description":"COMPRA A VISTA","debit":true},
			{"id":2,"description":"PAGAMENTO","debit":false}
		]
		`, rec.Body.String())
	}
}

func TestHandlerOperation_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationRepository := repository.NewMockOperations(ctrl)
	mockOperationRepository.EXPECT().FindAll(gomock.Any(), filter.OperationCollection{}).Return(nil, errors.New("err find all"))

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewOperation(OperationOpts{
		OperationRepository: mockOperationRepository,
	})

	assert.EqualError(t, h.FindAll(c), "code=400, message=err find all")
}
