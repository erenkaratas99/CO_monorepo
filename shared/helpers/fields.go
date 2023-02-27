package helpers

import (
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"net/http"
)

type CheckFields struct {
	ExactFields    []string `json:"exactFields"`
	MatchFields    []string `json:"matchFields"`
	SortableFields []string `json:"sortableFields"`
	Fields         []string `json:"fields"`
}

var checkFields = map[string]CheckFields{
	"fields": {
		SortableFields: []string{
			"order_date",
			"updated_at",
			"order_total"},
		MatchFields: []string{"order_item.name"},
		ExactFields: []string{
			"id",
			"customer_id",
			"payment_status",
			"shipment_status",
			"address.address_name",
			"address.address_line",
			"address.city",
			"address.country",
			"address.city_code",
			"order_item.id"},
		Fields: []string{
			"id",
			"customer_id",
			"order_date",
			"order_total",
			"payment_status",
			"shipment_status",
			"order_item",
			"address",
			"address.address_name",
			"address.address_line",
			"address.city",
			"address.country",
			"address.city_code",
			"order_item.id",
			"order_item.imageurl",
			"order_item.name",
			"updated_at"},
	},
}

func GetCheckFields() (*CheckFields, error) {
	checkFields, isExist := checkFields["fields"]
	if !isExist {
		return nil, customErrors.NewHTTPError(http.StatusInternalServerError,
			"FieldErr",
			"Fields could not have fetched correctly from the static list.")
	}
	return &checkFields, nil
}

func (cf *CheckFields) FieldsCompatibilityCheckRevised(fsf *FSF) error {
	var fieldsAppended []interface{}
	for _, fieldRequested := range fsf.ExactFilter.Key {
		fieldsAppended = append(fieldsAppended, fieldRequested)
	}
	for _, fieldRequested := range fsf.MatchFilter.Key {
		fieldsAppended = append(fieldsAppended, fieldRequested)
	}
	isField := IsSubset(fieldsAppended, cf.Fields)
	if !isField {
		return customErrors.NewHTTPError(400, "FieldErr",
			"Requested field for filtering not could have been found.")
	}
	isFieldExactComp := IsSubset(fsf.ExactFilter.Key, cf.ExactFields)
	if !isFieldExactComp {
		return customErrors.NewHTTPError(http.StatusBadRequest, "FieldErr",
			"Requested field for exact search is not compatible for that.")
	}
	isFieldMatchComp := IsSubset(fsf.MatchFilter.Key, cf.MatchFields)
	if !isFieldMatchComp {
		return customErrors.NewHTTPError(http.StatusBadRequest, "FieldErr",
			"Requested field for match filtering is not compatible for that.")
	}
	var valTempArr []interface{}
	if fsf.SortCond.Value != nil {
		valTempArr = append(valTempArr, fsf.SortCond.Value)
		isFieldSortable := ListRevisedIterator(valTempArr, cf.SortableFields)
		if !isFieldSortable {
			return customErrors.NewHTTPError(http.StatusBadRequest, "FieldErr",
				"Requested field for sorting is not compatible for that.")
		}
	}
	return nil
}
