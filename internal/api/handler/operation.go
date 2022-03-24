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
	OperationFindAllPath  = "/operations"
	OperationFindByIDPath = "/operations/:id"
	OperationCreatePath   = "/operations"
)

type (
	OperationOpts struct {
		OperationRepository repository.Operations
	}
	Operation struct {
		OperationOpts
	}
)

func NewOperation(opts OperationOpts) *Operation {
	return &Operation{opts}
}

func (o *Operation) Create(c echo.Context) error {
	ctx, span := jaeger.Span(c.Request().Context())
	defer span.End()

	request := &contract.OperationRequest{}
	if err := c.Bind(request); err != nil {
		c.Logger().Errorf("c.Bind failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := request.Validate(); err != nil {
		c.Logger().Errorf("request.Validate failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	typeOperation, _ := strconv.ParseBool(request.Debit)
	operationType, err := o.OperationRepository.Create(ctx, entity.Operation{
		Description: request.Description,
		Debit:       typeOperation,
	})
	if err != nil {
		c.Logger().Errorf("o.OperationTypeRepository.Create failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, operationType)
}

func (o *Operation) FindByID(c echo.Context) error {
	ctx, span := jaeger.Span(c.Request().Context())
	defer span.End()

	id, _ := strconv.Atoi(c.Param("id"))
	operationType, err := o.OperationRepository.FindByID(ctx, uint(id))
	if err != nil {
		c.Logger().Errorf("o.OperationTypeRepository.FindByID failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, operationType)
}

func (o *Operation) FindAll(c echo.Context) error {
	ctx, span := jaeger.Span(c.Request().Context())
	defer span.End()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))

	operationTypes, err := o.OperationRepository.FindAll(ctx, filter.OperationCollection{
		Page:        page,
		Size:        size,
		Description: c.QueryParam("description"),
		Debit:       c.QueryParam("debit"),
	})
	if err != nil {
		c.Logger().Errorf("o.OperationTypeRepository.FindAll failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, operationTypes)
}
