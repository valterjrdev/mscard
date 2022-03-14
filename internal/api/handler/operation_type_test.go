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

func TestHandlerOperationType_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.OperationType{
		ID:          1,
		Description: "COMPRA A VISTA",
		Negative:    true,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"description":"COMPRA A VISTA","negative":"true"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperationType(OperationTypeOpts{
		OperationTypeRepository: mockOperationTypeRepository,
	})

	if assert.NoError(t, h.Create(echo.New().NewContext(req, rec))) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"id":1,"description":"COMPRA A VISTA","negative":true}`, rec.Body.String())
	}
}

func TestHandlerOperationType_Create_Persist_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, repository.ErrOperationTypeCreate)

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"description":"COMPRA A VISTA","negative":"true"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperationType(OperationTypeOpts{
		OperationTypeRepository: mockOperationTypeRepository,
	})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=failed to create new operation type")
}

func TestHandlerOperationType_Create_BindRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperationType(OperationTypeOpts{})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=code=400, message=unexpected EOF, internal=unexpected EOF")
}

func TestHandlerOperationType_Create_Validate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, AccountCreatePath, strings.NewReader(`{"description":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := NewOperationType(OperationTypeOpts{})

	assert.EqualError(t, h.Create(echo.New().NewContext(req, rec)), "code=400, message=description: cannot be blank; negative: cannot be blank.")
}

func TestHandlerOperationType_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(&entity.OperationType{
		ID:          1,
		Description: "COMPRA A VISTA",
		Negative:    true,
	}, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetPath(AccountFindByIDPath)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewOperationType(OperationTypeOpts{
		OperationTypeRepository: mockOperationTypeRepository,
	})

	if assert.NoError(t, h.FindByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"id":1,"description":"COMPRA A VISTA","negative":true}`, rec.Body.String())
	}
}

func TestHandlerOperationType_FindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindByID(gomock.Any(), uint(1)).Return(nil, repository.ErrOperationTypeFindByID)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	c.SetPath(AccountFindByIDPath)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewOperationType(OperationTypeOpts{
		OperationTypeRepository: mockOperationTypeRepository,
	})

	assert.EqualError(t, h.FindByID(c), "code=400, message=failed fetch operation type")
}

func TestHandlerOperationType_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindAll(gomock.Any(), filter.OperationTypeCollection{}).Return([]*entity.OperationType{
		{
			ID:          1,
			Description: "COMPRA A VISTA",
			Negative:    true,
		},
		{
			ID:          2,
			Description: "PAGAMENTO",
			Negative:    false,
		},
	}, nil)

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewOperationType(OperationTypeOpts{
		OperationTypeRepository: mockOperationTypeRepository,
	})

	if assert.NoError(t, h.FindAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `
		[
			{"id":1,"description":"COMPRA A VISTA","negative":true},
			{"id":2,"description":"PAGAMENTO","negative":false}
		]
		`, rec.Body.String())
	}
}

func TestHandlerOperationType_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperationTypeRepository := repository.NewMockOperationTypes(ctrl)
	mockOperationTypeRepository.EXPECT().FindAll(gomock.Any(), filter.OperationTypeCollection{}).Return(nil, errors.New("err find all"))

	server := echo.New()
	req := httptest.NewRequest(http.MethodPost, AccountFindByIDPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := server.NewContext(req, rec)
	h := NewOperationType(OperationTypeOpts{
		OperationTypeRepository: mockOperationTypeRepository,
	})

	assert.EqualError(t, h.FindAll(c), "code=400, message=err find all")
}
