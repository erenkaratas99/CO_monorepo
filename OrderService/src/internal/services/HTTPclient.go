package services

import (
	sharedentities "CustomerOrderMonoRepo/shared/entities"
)

func (s *Service) orderClient(customerid, corID string) (*sharedentities.Address, error) {
	address := sharedentities.Address{}
	URI := "http://localhost:8000/customers/address/" + customerid
	err := s.client.DoGetRequest(URI, corID, &address)
	if err != nil {
		return nil, err
	}
	return &address, nil
}
