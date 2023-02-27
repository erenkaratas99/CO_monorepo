package helpers

import sharedentities "CustomerOrderMonoRepo/shared/entities"

func NewResponseModel(totalObj *int, batchLen *int, data interface{}) *sharedentities.ResponseModel {
	model := sharedentities.ResponseModel{
		RespObjectCount:  batchLen,
		TotalObjectCount: totalObj,
		Data:             &data,
	}
	return &model
}
