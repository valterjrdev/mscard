package handler

import (
	"github.com/labstack/echo/v4"
	"ms/card/pkg/contract"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/filter"
	"ms/card/pkg/persistence/repository"
	"net/http"
	"strconv"
)

const (
	OperationTypeFindAllPath  = "/operation-types"
	OperationTypeFindByIDPath = "/operation-types/:id"
	OperationTypeCreatePath   = "/operation-types"
)

type (
	OperationTypeOpts struct {
		OperationTypeRepository repository.OperationTypes
	}
	OperationType struct {
		OperationTypeOpts
	}
)

func NewOperationType(opts OperationTypeOpts) *OperationType {
	return &OperationType{opts}
}

func (o *OperationType) Create(c echo.Context) error {
	request := &contract.OperationTypeRequest{}
	if err := c.Bind(request); err != nil {
		c.Logger().Errorf("c.Bind failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := request.Validate(); err != nil {
		c.Logger().Errorf("request.Validate failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	negative, _ := strconv.ParseBool(request.Negative)
	operationType, err := o.OperationTypeRepository.Create(c.Request().Context(), entity.OperationType{
		Description: request.Description,
		Negative:    negative,
	})
	if err != nil {
		c.Logger().Errorf("o.OperationTypeRepository.Create failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, operationType)
}

func (o *OperationType) FindByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	operationType, err := o.OperationTypeRepository.FindByID(c.Request().Context(), uint(id))
	if err != nil {
		c.Logger().Errorf("o.OperationTypeRepository.FindByID failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, operationType)
}

func (o *OperationType) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))

	operationTypes, err := o.OperationTypeRepository.FindAll(c.Request().Context(), filter.OperationTypeCollection{
		Page:        page,
		Size:        size,
		Description: c.QueryParam("description"),
		Negative:    c.QueryParam("negative"),
	})
	if err != nil {
		c.Logger().Errorf("o.OperationTypeRepository.FindAll failed with %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, operationTypes)
}
