package handlers

import (
	"CustomerOrderMonoRepo/CustomerService/src/internal/entities"
	"CustomerOrderMonoRepo/CustomerService/src/internal/repositories"
	"CustomerOrderMonoRepo/CustomerService/src/internal/services"
	_ "CustomerOrderMonoRepo/docs"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericHandler"
	"github.com/erenkaratas99/COApiCore/pkg"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"github.com/erenkaratas99/COApiCore/pkg/helper"
	"github.com/erenkaratas99/COApiCore/shared/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type Handler struct {
	repo           *repositories.Repository
	service        *services.Service
	echo           *echo.Echo
	genericHandler *genericHandler.Handler
}

func NewHandler(repo *repositories.Repository, service *services.Service, echo *echo.Echo, genericHandler *genericHandler.Handler) *Handler { //
	return &Handler{repo: repo, service: service, echo: echo, genericHandler: genericHandler} //
}

func (h *Handler) InitEndpoints() {
	e := h.echo
	g := e.Group("/customers")
	g.POST("/", h.CreateCustomer)
	g.GET("/:customerid", h.GetbyID)
	g.GET("/", h.GetAll)
	g.PUT("/:customerid", h.UpdateCustomer)
	g.DELETE("/:customerid", h.DeleteCustomer)
	g.GET("/address/:customerid", h.GetAddress)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	g.GET("/generic", h.genericHandler.GenericGet)
}

// CreateCustomer godoc
// @Summary      Creates a new customer
// @Description  Creates a new customer with the given data
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customerReq   body    entities.CustomerRequestModel   true  "Customer Request Model"
// @Success      201
// @Router       /customers/ [post]
func (h *Handler) CreateCustomer(c echo.Context) error {
	customerReq := entities.CustomerRequestModel{}
	err := c.Bind(&customerReq)
	if err != nil {
		return customErrors.BindErr
	}
	err = c.Validate(customerReq)
	if err != nil {
		return err
	}
	err, insertedID := h.service.CreateCustomerService(&customerReq)
	if err != nil {
		return err
	}
	srm := types.GetSRM(*insertedID)
	return c.JSON(http.StatusCreated, srm)
}

// GetbyID godoc
// @Summary      Gets a single customer
// @Description  Gets a single customer with the given ID
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customerid    path    string     true    "Customer ID"
// @Success      200  {object}  entities.Customer
// @Router       /customers/:customerid [get]
func (h *Handler) GetbyID(c echo.Context) error {
	id := c.Param("customerid")
	_, err := uuid.Parse(id)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	customerResp, err := h.service.GetSingleCustomerService(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, *customerResp)
}

// GetAll godoc
// @Summary      Gets all customers
// @Description  Gets all customers with pagination support
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        limit    query    string    false   "Limit"
// @Param        offset   query    string    false   "Offset"
// @Success      200  {object}  []entities.Customer
// @Router       /customers/ [get]
func (h *Handler) GetAll(c echo.Context) error {
	l := c.QueryParam("limit")
	o := c.QueryParam("offset")
	limit, offset := pkg.LimitOffsetValidation(l, o)
	customers, err := h.service.GetAllCustomersService(limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customers)
}

// UpdateCustomer godoc
// @Summary      Updates an existing customer
// @Description  Updates an existing customer with the given data
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customerid    path    string                     true    "Customer ID"
// @Param        customerReq   body    entities.CustomerRequestModel  true    "Customer Request Model"
// @Success      201
// @Router       /customers/:customerid [put]
func (h *Handler) UpdateCustomer(c echo.Context) error {
	id := c.Param("customerid")
	_, err := uuid.Parse(id)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	customerReq := entities.CustomerRequestModel{}
	err = c.Bind(&customerReq)
	if err != nil {
		return err
	}
	err = c.Validate(customerReq)
	if err != nil {
		return err
	}
	err, upsertedID := h.service.UpdateCustomerService(id, &customerReq)
	if err != nil {
		return err
	}
	srm := types.GetSRM(*upsertedID)
	return c.JSON(http.StatusCreated, srm)
}

// DeleteCustomer godoc
// @Summary      Deletes an existing customer
// @Description  Deletes an existing customer with the given ID
// @Description  Checks the order service whether if the user has non-delivered orders
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customerid    path    string     true    "Customer ID"
// @Success      204  {string}  string
// @Router       /customers/:customerid [delete]
func (h *Handler) DeleteCustomer(c echo.Context) error {
	id := c.Param("customerid")
	_, err := uuid.Parse(id)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	corID := helper.GetCorrelationID(c)
	err = h.service.DeleteCustomerService(id, corID)
	if err != nil {
		return err
	}
	return nil
}

// GetAddress godoc
// @Summary      Gets the address of a customer
// @Description  Gets the address of a customer with the given ID
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customerid    path    string     true    "Customer ID"
// @Success      200 {object} sharedentities.Address
// @Router       /customers/address/:customerid [get]
func (h *Handler) GetAddress(c echo.Context) error {
	customerid := c.Param("customerid")
	_, err := uuid.Parse(customerid)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	address, err := h.service.GetAddressService(customerid)
	if err != nil {
		return err
	}
	if address != nil {
		return c.JSON(http.StatusOK, address)
	}
	return c.JSON(http.StatusOK, nil)
}
