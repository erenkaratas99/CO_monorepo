package entities

type OrderResponseModel struct {
	Id             string `json:"id"bson:"_id"`
	CustomerId     string `json:"customer_id,omitempty"bson:"customer_id"`
	OrderDate      string `json:"order_date,omitempty"bson:"order_date"`
	UpdatedAt      string `json:"updated_at,omitempty"bson:"updated_at"`
	PaymentStatus  string `json:"payment_status,omitempty"bson:"payment_status"`
	ShipmentStatus string `json:"shipment_status,omitempty"bson:"shipment_status"`
}

type OrderRequestModel struct {
	OrderItem      *Product `json:"order_item"bson:"order_item"`
	OrderTotal     int      `json:"order_total"bson:"order_total"`
	ShipmentStatus string   `json:"shipment_status"bson:"shipment_status"`
	PaymentStatus  string   `json:"payment_status"bson:"payment_status"`
}
