package handler

import (
	"github.com/labstack/echo/v4"
	"ms/card/pkg/contract"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
	"ms/card/pkg/persistence/repository"
	"ms/card/pkg/telemetry/jaeger"
	"net/http"
	"strconv"
)

const (
	AccountFindAllPath  = "/accounts"
	AccountFindByIDPath = "/accounts/:id"
	AccountCreatePath   = "/accounts"
)

type (
	AccountOpts struct {
		AccountRepository repository.Accounts
	}
	Account struct {
		AccountOpts
	}
)

func NewAccount(opts AccountOpts) *Account {
	return &Account{opts}
}

func (a *Account) Create(c echo.Context) error {
	ctx, span := jaeger.Span(c.Request().Context())
	defer span.End()

	request := &contract.AccountRequest{}
	if err := c.Bind(request); err != nil {
		c.Logger().Errorf("c.Bind failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := request.Validate(); err != nil {
		c.Logger().Errorf("request.Validate failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	account, err := a.AccountRepository.Create(ctx, entity.Account{
		Document: request.Document,
		Limit:    request.Limit,
	})
	if err != nil {
		c.Logger().Errorf("a.AccountRepository.Create failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, account)
}

func (a *Account) FindByID(c echo.Context) error {
	ctx, span := jaeger.Span(c.Request().Context())
	defer span.End()

	id, _ := strconv.Atoi(c.Param("id"))
	account, err := a.AccountRepository.FindByID(ctx, uint(id))
	if err != nil {
		c.Logger().Errorf("a.AccountRepository.FindByID failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

func (a *Account) FindAll(c echo.Context) error {
	ctx, span := jaeger.Span(c.Request().Context())
	defer span.End()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	accounts, err := a.AccountRepository.FindAll(ctx, filter.AccountCollection{
		Page:     page,
		Size:     size,
		Document: c.QueryParam("document_number"),
	})
	if err != nil {
		c.Logger().Errorf("a.AccountRepository.FindAll failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, accounts)
}
