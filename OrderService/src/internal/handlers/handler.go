package handlers

import (
	//_ "CustomerOrderMonoRepo/OrderService/src/docs"
	"CustomerOrderMonoRepo/OrderService/src/internal/entities"
	"CustomerOrderMonoRepo/OrderService/src/internal/helpers"
	"CustomerOrderMonoRepo/OrderService/src/internal/repositories"
	"CustomerOrderMonoRepo/OrderService/src/internal/services"
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

func NewHandler(repo *repositories.Repository, service *services.Service, echo *echo.Echo, ghandler *genericHandler.Handler) *Handler {
	return &Handler{repo: repo, service: service, echo: echo, genericHandler: ghandler}
}

func (h *Handler) InitEndpoints() {
	e := h.echo
	g := e.Group("/orders")
	g.POST("/:customerid", h.CreateOrder)
	g.GET("/:orderid", h.GetByID)
	g.GET("/", h.GetAll)
	g.DELETE("/:orderid", h.DeleteOrder)
	g.PUT("/orderid", h.UpdateOrder)
	g.GET("/orderof/:customerid/", h.GetCustomerOrders)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	g.GET("/generic", h.genericHandler.GenericGet)
}

// CreateOrder godoc
// @Summary      It creates an order
// @Description  It creates an order that comes with req. Body as JSON.
// @Description  has validation
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        entities.OrderRequestModel  body   entities.OrderRequestModel  true  "Order fields"
// @Success      200      {string}  string         "OK"
// @Failure      500      {error}  error         "binding error"
// @Failure      400      {error}  error         "bad request"
// @Router       /orders/:customerid [post]
func (h *Handler) CreateOrder(c echo.Context) error {
	orderReq := entities.OrderRequestModel{}
	customerid := c.Param("customerid")
	_, err := uuid.Parse(customerid)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	corId := helper.GetCorrelationID(c)
	err = c.Bind(&orderReq)
	if err != nil {
		return customErrors.BindErr
	}
	err = c.Validate(orderReq)
	if err != nil {
		return err
	}
	err = helpers.ShipStatProdIDCheck(&orderReq)
	if err != nil {
		return err
	}
	err, insertedID := h.service.CreateOrderService(customerid, corId, &orderReq)
	if err != nil {
		return err
	}
	srm := types.GetSRM(*insertedID)
	return c.JSON(http.StatusCreated, srm)
}

// GetByID godoc
// @Summary      It serves "an" order
// @Description  It gets an order due to its UUID
// @Tags         orders
// @Accept 	     json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  entities.Order
// @Router       /orders/:orderid [get]
func (h *Handler) GetByID(c echo.Context) error {
	orderid := c.Param("orderid")
	_, err := uuid.Parse(orderid)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	orderResp, err := h.service.GetSingleOrderService(orderid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, *orderResp)
}

// GetAll godoc
// @Summary      It serves all orders
// @Description  It gets all order list due to limit offset values
// @Tags         orders
// @Accept 	     json
// @Produce      json
// @Param        limit   query  string  false "limit"
// @Param        offset   query  string  false "offset"
// @Success      200  {object}  entities.Order
// @Failure      500  {error}   error "internal server error"
// @Router       /orders [get]
func (h *Handler) GetAll(c echo.Context) error {
	l := c.QueryParam("limit")
	o := c.QueryParam("offset")
	limit, offset := pkg.LimitOffsetValidation(l, o)
	orders, err := h.service.GetAllOrdersService(limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orders)
}

// DeleteOrder godoc
// @Summary      It deletes "an" order
// @Description  It deletes an order due to its UUID
// @Tags         orders
// @Accept 	     json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {string}  string "1 order has been deleted"
// @Failure      400  {error}  error "bad request"
// @Failure      404  {error}  error "Given ID param does not match any order."
// @Router       /orders/:orderid [delete]
func (h *Handler) DeleteOrder(c echo.Context) error {
	orderid := c.Param("orderid")
	_, err := uuid.Parse(orderid)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	err = h.service.DeleteOrderService(orderid)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrder godoc
// @Summary      It updates an order
// @Description  It updates an order that comes with req. Body as JSON.
// @Description  has validation
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        entities.OrderRequestModel  body   entities.OrderRequestModel  true  "order fields"
// @Param        ID  path   string  true  "Order UUID"
// @Success      201
// @Failure      500      {error}  error         "binding error"
// @Failure      400      {error}  error         "bad request"
// @Router       /orders/:orderid [put]
func (h *Handler) UpdateOrder(c echo.Context) error {
	orderid := c.Param("orderid")
	_, err := uuid.Parse(orderid)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	orderReq := entities.OrderRequestModel{}
	err = c.Bind(&orderReq)
	if err != nil {
		return customErrors.BindErr
	}
	err = helpers.ShipStatProdIDCheck(&orderReq)
	if err != nil {
		return err
	}
	err = c.Validate(orderReq)
	if err != nil {
		return err
	}
	err, upsertedID := h.service.UpdateOrderService(orderid, &orderReq)
	if err != nil {
		return err
	}
	srm := types.GetSRM(*upsertedID)
	return c.JSON(http.StatusCreated, srm)
}

// GetCustomerOrders godoc
// @Summary      Gets the orders of a customer
// @Description  Gets the orders of a customer with the given ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        customerid    path    string     true    "Customer ID"
// @Success      200  {object}  sharedentities.Address
// @Router       /orders/orderof/:customerid [get]
func (h *Handler) GetCustomerOrders(c echo.Context) error {
	l := c.QueryParam("limit")
	o := c.QueryParam("offset")
	customerId := c.Param("customerid")
	_, err := uuid.Parse(customerId)
	if err != nil {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	limit, offset := pkg.LimitOffsetValidation(l, o)
	orders, err := h.service.GetCustomerOrdersService(customerId, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orders)
}
