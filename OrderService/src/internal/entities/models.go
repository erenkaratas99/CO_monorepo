package entities

import sharedentities "CustomerOrderMonoRepo/shared/entities"

type Order struct {
	Id             string                  `json:"id" bson:"_id"`
	CustomerID     string                  `json:"customer_id" bson:"customer_id"`
	OrderDate      string                  `json:"order_date"bson:"order_date"`
	OrderTotal     int                     `json:"order_total"bson:"order_total"`
	PaymentStatus  string                  `json:"payment_status"bson:"payment_status"`
	ShipmentStatus string                  `json:"shipment_status"bson:"shipment_status"`
	OrderItem      *Product                `json:"order_item"bson:"order_item"`
	Address        *sharedentities.Address `json:"address"bson:"address"`
	UpdatedAt      string                  `json:"updated_at"bson:"updated_at"`
}

type Product struct {
	Id     string `json:"id" bson:"_id"`
	ImgUrl string `json:"imageurl" bson:"imageurl"`
	Name   string `json:"name"bson:"name"`
}
