package sharedentities

// sharedentities.Address model info
// @Description a nested struct for a field of entities.Customer and entities.Order
type Address struct {
	CustomerId  string `json:"customer_id"bson:"customer_id"`
	AddressName string `json:"address_name"bson:"address_name"`
	AddressLine string `json:"address_line"bson:"address_line"`
	City        string `json:"city"bson:"city"`
	Country     string `json:"country"bson:"country"`
	CityCode    int    `json:"city_code"bson:"city_code"`
}

type JsonMap struct {
	ExactFilters map[string][]interface{} `json:"exactFilters"`
	Match        map[string][]interface{} `json:"match"`
	Fields       map[string][]interface{} `json:"fields"`
	SortCond     map[string]interface{}   `json:"sortCond"`
}

// sharedentities.ResponseModel model info
// @Description a struct for response with total object count and response object count fields
type ResponseModel struct {
	RespObjectCount  *int         `json:"resp_object_count"`
	TotalObjectCount *int         `json:"total_object_count"`
	Data             *interface{} `json:"data"`
}
