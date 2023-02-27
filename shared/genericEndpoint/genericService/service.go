package genericService

import (
	sharedentities "CustomerOrderMonoRepo/shared/entities"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericRepository"
	"CustomerOrderMonoRepo/shared/helpers"
)

type Service struct {
	repo *genericRepository.Repository
}

func NewGenericService(r *genericRepository.Repository) *Service {
	return &Service{repo: r}
}
func (s *Service) GenericGetService(jsonMap sharedentities.JsonMap) (*sharedentities.ResponseModel, error) { // []*entities.Order,
	fsf := helpers.NewFSF(jsonMap)
	cf, err := helpers.GetCheckFields()
	if err != nil {
		return nil, err
	}
	err = cf.FieldsCompatibilityCheckRevised(fsf)
	if err != nil {
		return nil, err
	}
	order, totalCount, err := s.repo.GenericGetRepo(fsf, false)
	if err != nil {
		return nil, err
	}
	length := len(order)
	resp := helpers.NewResponseModel(totalCount, &length, order)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
