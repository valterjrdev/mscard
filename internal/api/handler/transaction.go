package handler

import (
	"github.com/labstack/echo/v4"
	"ms/card/pkg/contract"
	"ms/card/pkg/persistence/filter"
	"ms/card/pkg/persistence/repository"
	"ms/card/pkg/service"
	"net/http"
	"strconv"
)

const (
	TransactionFindAllPath = "/transactions"
	TransactionCreatePath  = "/transactions"
)

type (
	TransactionOpts struct {
		TransactionService    service.Transactions
		TransactionRepository repository.Transactions
	}
	Transaction struct {
		TransactionOpts
	}
)

func NewTransaction(opts TransactionOpts) *Transaction {
	return &Transaction{opts}
}

func (t *Transaction) Create(c echo.Context) error {
	request := &contract.TransactionRequest{}
	if err := c.Bind(request); err != nil {
		c.Logger().Errorf("c.Bind failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transaction, err := t.TransactionService.Create(c.Request().Context(), request)
	if err != nil {
		c.Logger().Errorf("t.TransactionService.Create failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (t *Transaction) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	eventDateStart := c.QueryParam("eventDateStart")
	eventDateEnd := c.QueryParam("eventDateEnd")

	collection, err := t.TransactionRepository.FindAll(c.Request().Context(), filter.TransactionCollection{
		Page:           page,
		Size:           size,
		Account:        c.QueryParam("account_id"),
		Operation:      c.QueryParam("operation_id"),
		EventDateStart: eventDateStart,
		EventDateEnd:   eventDateEnd,
	})
	if err != nil {
		c.Logger().Errorf("t.TransactionRepository.FindAll failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, collection)
}
