package entities

import (
	sharedentities "CustomerOrderMonoRepo/shared/entities"
)

type CustomerRequestModel struct {
	FirstName string                  `json:"first_name"bson:"first_name"`
	LastName  string                  `json:"last_name"bson:"last_name"`
	Email     string                  `json:"email"bson:"email"`
	Phone     string                  `json:"phone"bson:"phone"`
	Address   *sharedentities.Address `json:"address"bson:"address"`
}

type CustomerResponseModel struct {
	FirstName string `json:"first_name,omitempty"bson:"first_name"`
	LastName  string `json:"last_name,omitempty"bson:"last_name"`
	Email     string `json:"email,omitempty"bson:"email"`
	CreatedAt string `json:"created_at,omitempty"bson:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"bson:"updated_at"`
}
type CustomerAddressResponseModel struct {
	*sharedentities.Address
}
