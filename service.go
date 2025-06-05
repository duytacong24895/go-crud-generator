package crud_generator

import (
	"fmt"

	"github.com/duytacong24895/go-curd-generator/core"
	"gorm.io/gorm"
)

type IService interface {
	Create(model *core.Model, inputData *map[string]any) (any, error)
	GetByID(model *core.Model, id string) (any, error)
	GetList(model *core.Model, inputData *GetListQueryParams) ([]*map[string]any, int64, error)
	Update(model *core.Model, inputData *map[string]any, id string) (*map[string]any, error)
	Delete(model *core.Model, id string) error
}
type service struct {
	repository Repository
}

func NewService(repository Repository) IService {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(model *core.Model, inputData *map[string]any) (any, error) {
	entity, err := s.repository.Create(model, inputData)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) GetByID(model *core.Model, id string) (any, error) {
	entity, err := s.repository.GetByID(model, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) GetList(model *core.Model, inputData *GetListQueryParams) ([]*map[string]any, int64, error) {
	entities, total, err := s.repository.GetList(model, inputData.Page, inputData.PageSize,
		inputData.Filter, inputData.OrderBy)
	if err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}

func (s *service) Update(model *core.Model, inputData *map[string]any, id string) (*map[string]any, error) {
	_, err := s.repository.GetByID(model, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}

	entities, err := s.repository.Update(model, inputData, id)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (s *service) Delete(model *core.Model, id string) error {
	return s.repository.Delete(model, id)
}
