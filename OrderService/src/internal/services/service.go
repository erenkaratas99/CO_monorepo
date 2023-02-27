package services

import (
	"CustomerOrderMonoRepo/OrderService/src/internal/entities"
	"CustomerOrderMonoRepo/OrderService/src/internal/repositories"
	sharedentities "CustomerOrderMonoRepo/shared/entities"
	"CustomerOrderMonoRepo/shared/helpers"
	"github.com/erenkaratas99/COApiCore/pkg"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
)

type Service struct {
	repo   *repositories.Repository
	client *pkg.RestClient
}

func NewService(r *repositories.Repository, client *pkg.RestClient) *Service {
	return &Service{repo: r, client: client}
}

func (s *Service) CreateOrderService(customerid, corID string, orderReq *entities.OrderRequestModel) (error, *string) {
	address, err := s.orderClient(customerid, corID)
	if err != nil {
		return err, nil
	}
	order := &entities.Order{
		CustomerID:     customerid,
		OrderTotal:     orderReq.OrderTotal,
		PaymentStatus:  orderReq.PaymentStatus,
		ShipmentStatus: orderReq.ShipmentStatus,
		Address:        address,
		OrderItem:      orderReq.OrderItem,
	}
	insertedID, err := s.repo.CreateOrder(order)
	if err != nil {
		return err, nil
	}
	return nil, insertedID
}

func (s *Service) GetSingleOrderService(orderid string) (*sharedentities.ResponseModel, error) {
	order, totalCount, err := s.repo.GetSingleOrder(orderid, true)
	if err != nil {
		return nil, err
	}
	orderResp := &entities.OrderResponseModel{
		Id:             orderid,
		CustomerId:     order.CustomerID,
		OrderDate:      order.OrderDate,
		UpdatedAt:      order.UpdatedAt,
		PaymentStatus:  order.PaymentStatus,
		ShipmentStatus: order.ShipmentStatus,
	}
	length := 1
	resp := helpers.NewResponseModel(totalCount, &length, orderResp)
	return resp, nil
}

func (s *Service) GetAllOrdersService(l, o int64) (*sharedentities.ResponseModel, error) {
	orders, totalCount, err := s.repo.GetAllOrders(l, o, true)
	if err != nil {
		return nil, err
	}
	length := len(orders)
	resp := helpers.NewResponseModel(totalCount, &length, orders)
	return resp, nil
}

func (s *Service) DeleteOrderService(orderid string) error {
	err := s.repo.DeleteOrder(orderid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateOrderService(orderid string, orderReq *entities.OrderRequestModel) (error, *string) {
	exist, err := helpers.ObjectExists(s.repo.Collection, orderid)
	if err != nil {
		return err, nil
	}
	if exist {
		order := &entities.Order{
			OrderTotal:     orderReq.OrderTotal,
			PaymentStatus:  orderReq.PaymentStatus,
			ShipmentStatus: orderReq.ShipmentStatus,
			OrderItem:      orderReq.OrderItem,
		}
		err, upsertedID := s.repo.UpdateOrder(orderid, order)
		if err != nil {
			return err, nil
		}
		return nil, upsertedID
	}
	return customErrors.DocNotFound, nil
}

func (s *Service) GetCustomerOrdersService(customerID string, l, o int64) ([]*entities.OrderResponseModel, error) {
	orders, err := s.repo.GetCustomerOrders(customerID, l, o)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
