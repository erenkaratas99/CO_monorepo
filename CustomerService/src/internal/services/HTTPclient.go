package services

import (
	"github.com/erenkaratas99/COApiCore/pkg"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"net/http"
)

type customerClientResponse struct {
	Id             string `json:"id"`
	ShipmentStatus string `json:"shipment_status"`
}

func (s *Service) customerClient(id, corID string) ([]customerClientResponse, error) {
	orderList := []customerClientResponse{}
	limitOffsetParams := pkg.BuildLimitOffsetParams(25, 0)
	URI := "http://localhost:8001/orders/orderof/" + id + limitOffsetParams
	err := s.client.DoGetRequest(URI, corID, &orderList)
	if err != nil {
		return nil, err
	}
	return orderList, nil
}

func ValidateForDeletion(orderlist []customerClientResponse) error {
	if len(orderlist) == 0 {
		return nil
	}
	for _, order := range orderlist {
		if order.ShipmentStatus != "delivered" || order.ShipmentStatus == "cancelled" {
			return customErrors.NewHTTPError(http.StatusConflict,
				"CustomerErr", "Customer has orders on progress.")
		}
	}
	return nil
}
