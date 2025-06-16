package repositories

import (
	"time"

	"github.com/duytacong24895/go-crud-generator/core"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepository interface {
	GetList(model *core.Model, page int, pageSize int,
		filter core.IFilter, order_by string) ([]*map[string]any, int64, error)
	Create(model *core.Model, inputData *map[string]any) (*map[string]any, error)
	GetByID(model *core.Model, id string) (*map[string]any, error)
	Update(model *core.Model, inputData *map[string]any, id string) (*map[string]any, error)
	Delete(model *core.Model, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(model *core.Model, inputData *map[string]any) (*map[string]any, error) {
	now := time.Now()
	if model.Meta.UpdatedAtField != nil {
		(*inputData)[model.Meta.UpdatedAtField.Name] = now
	}

	if model.Meta.CreatedAtField != nil {
		(*inputData)[model.Meta.CreatedAtField.Name] = now
	}

	if err := r.db.Clauses(clause.Returning{}).Model(&model.Ref).Create(inputData).Error; err != nil {
		return nil, err
	}
	return inputData, nil
}

func (r *repository) GetByID(model *core.Model, id string) (*map[string]any, error) {
	var entity = make(map[string]any)
	statement := r.db.Model(&model.Ref).Where("id = ?", id)
	if model.Meta.SoftDeletedField != nil {
		statement = statement.Where(model.Meta.SoftDeletedField.DBName + " IS NULL")
	}
	if err := statement.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) GetList(model *core.Model, page int, pageSize int,
	filter core.IFilter, order_by string) ([]*map[string]any, int64, error) {
	var entities = make([]map[string]any, 0)
	var queryStatement *gorm.DB
	if filter.IsEmpty() {
		queryStatement = r.db.Model(&model.Ref)
	} else {
		var err error
		queryStatement, err = filter.BuildQuery(r.db)
		if err != nil {
			return nil, 0, err
		}
		queryStatement = queryStatement.Model(&model.Ref)
	}

	if order_by != "" {
		queryStatement = queryStatement.Order(order_by)
	}

	if model.Meta.SoftDeletedField != nil {
		queryStatement = queryStatement.Where(model.Meta.SoftDeletedField.DBName + " IS NULL")
	}

	var total int64
	if err := queryStatement.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := queryStatement.Debug().Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	var result []*map[string]any
	for i := range entities {
		result = append(result, &entities[i])
	}
	return result, total, nil
}

func (r *repository) Update(model *core.Model, inputData *map[string]any, id string) (*map[string]any, error) {
	if model.Meta.UpdatedAtField != nil {
		// Soft delete
		(*inputData)[model.Meta.UpdatedAtField.Name] = time.Now()
	}

	if err := r.db.Model(&model.Ref).Where("id = ?", id).Updates(&inputData).Error; err != nil {
		return nil, err
	}
	return inputData, nil
}

func (r *repository) Delete(model *core.Model, id string) error {
	if model.Meta.SoftDeletedField != nil {
		// Soft delete
		return r.db.Model(&model.Ref).Where("id = ?", id).
			Update(model.Meta.SoftDeletedField.Name, time.Now()).Error
	}
	return r.db.Model(&model.Ref).Where("id = ?", id).Delete(model.Ref).Error
}
