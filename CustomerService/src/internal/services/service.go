package services

import (
	"CustomerOrderMonoRepo/CustomerService/src/internal/entities"
	"CustomerOrderMonoRepo/CustomerService/src/internal/repositories"
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

func (s *Service) CreateCustomerService(customerReq *entities.CustomerRequestModel) (error, *string) {
	insertedID, err := s.repo.InsertCustomer(customerReq)
	if err != nil {
		return err, nil
	}
	return nil, insertedID
}

func (s *Service) GetSingleCustomerService(id string) (*sharedentities.ResponseModel, error) {
	customers, totalCount, err := s.repo.GetSingleCustomer(id, true)
	len := 1
	customerResp := helpers.NewResponseModel(totalCount, &len, customers)
	if err != nil {
		return nil, err
	}
	return customerResp, nil

}

func (s *Service) GetAllCustomersService(l, o int64) (*sharedentities.ResponseModel, error) {
	customers, totalCount, err := s.repo.GetAllCustomers(l, o, true)
	if err != nil {
		return nil, err
	}
	length := len(customers)
	resp := helpers.NewResponseModel(totalCount, &length, customers)
	return resp, nil
}

func (s *Service) UpdateCustomerService(id string, customerReq *entities.CustomerRequestModel) (error, *string) {
	exist, err := helpers.ObjectExists(s.repo.Collection, id)
	if err != nil {
		return err, nil
	}
	if exist {
		customer := &entities.Customer{
			Id:        id,
			FirstName: customerReq.FirstName,
			LastName:  customerReq.LastName,
			Email:     customerReq.Email,
			Phone:     customerReq.Phone,
			Address:   customerReq.Address,
		}
		err, upsertedID := s.repo.UpdateCustomer(customer)
		if err != nil {
			return err, nil
		}
		return nil, upsertedID
	}
	return customErrors.DocNotFound, nil
}

func (s *Service) DeleteCustomerService(id, corID string) error {
	exist, err := helpers.ObjectExists(s.repo.Collection, id)
	if err != nil {
		return err
	}
	if exist {
		orderlist, err := s.customerClient(id, corID)
		if err != nil {
			return err
		}
		err = ValidateForDeletion(orderlist)
		if err != nil {
			return err
		}
		err = s.repo.DeleteCustomer(id)
		if err != nil {
			return err
		}
		return nil
	}
	return customErrors.DocNotFound
}

func (s *Service) GetAddressService(id string) (*entities.CustomerAddressResponseModel, error) {
	address, err := s.repo.GetAddressOfCustomer(id)
	if err != nil {
		return nil, err
	}
	return address, nil
}
