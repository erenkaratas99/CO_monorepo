package entities

import sharedentities "CustomerOrderMonoRepo/shared/entities"

type Customer struct {
	Id        string                  `json:"id" bson:"_id"`
	FirstName string                  `json:"first_name"bson:"first_name"`
	LastName  string                  `json:"last_name"bson:"last_name"`
	Email     string                  `json:"email"bson:"email"`
	Phone     string                  `json:"phone"bson:"phone"`
	Address   *sharedentities.Address `json:"address"bson:"address"`
	CreatedAt string                  `json:"created_at"bson:"created_at"`
	UpdatedAt string                  `json:"updated_at"bson:"updated_at"`
}
