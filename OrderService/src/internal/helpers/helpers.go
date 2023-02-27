package helpers

import (
	"CustomerOrderMonoRepo/OrderService/src/internal/entities"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"net/http"
)

func ShipStatProdIDCheck(requestBody *entities.OrderRequestModel) error {
	validStatus := map[string]bool{
		"delivery":  true,
		"delivered": true,
		"shipment":  true,
		"cancelled": true,
	}
	if !validStatus[requestBody.ShipmentStatus] {
		return customErrors.NewHTTPError(http.StatusBadRequest, "ShipmentStatusErr",
			"shipment_status must be one of : ['delivery', 'delivered', 'shipment', 'cancelled']")
	}
	productIdList := map[string]bool{
		"1139": true,
		"1971": true,
		"1999": true,
		"1938": true,
		"1923": true,
		"1881": true,
	}
	if !productIdList[requestBody.OrderItem.Id] {
		return customErrors.NewHTTPError(http.StatusBadRequest, "ProductIDErr",
			"OrderItem id (product id) is not valid")
	}
	return nil
}
