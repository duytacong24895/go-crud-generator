package core

import (
	"reflect"
	"slices"
	"strings"

	constants "github.com/duytacong24895/go-curd-generator/const"
)

type Model struct {
	Name string `json:"name"`
	Ref  any
	Meta *MetaModel
}

type MetaModel struct {
	SoftDeletedField *ModelField
	CreatedAtField   *ModelField
	UpdatedAtField   *ModelField
}

type ModelField struct {
	Name   string `json:"name"` // Name of the field in the struct
	DBName string `json:"db_name"`
}

func NewModel(ref any) *Model {
	core := &Core{}
	return &Model{
		Name: core.ExactModelName(ref),
		Ref:  ref,
		Meta: NewMetaModel(ref),
	}
}

func NewMetaModel(ref any) *MetaModel {
	/*
		This function extracts metadata from the model reference.
		It looks for specific struct tags to identify fields related to soft deletion,
		creation, and update timestamps.

		If you are using from gorm.Model in your struct, You don't need to set these tags.
	*/
	gormSchema := Core{}.ExactSchemaGorm(ref)
	numField := reflect.TypeOf(ref).Elem().NumField()
	var softDeletedField, createdAtField, updatedAtField *ModelField
	for i := 0; i < numField; i++ {
		field := reflect.TypeOf(ref).Elem().Field(i)
		tags := field.Tag.Get(constants.FieldTagKey)
		if tags == "" {
			continue
		}
		arrTags := strings.Split(tags, constants.SepOfTags)
		if slices.Contains(arrTags, constants.SoftDeleteFieldTagName) {
			softDeletedField = &ModelField{
				Name:   field.Name,
				DBName: gormSchema[field.Name],
			}
		} else if slices.Contains(arrTags, constants.CreateTimeFieldTagName) {
			createdAtField = &ModelField{
				Name:   field.Name,
				DBName: gormSchema[field.Name],
			}
		} else if slices.Contains(arrTags, constants.UpdateTimeFieldTagName) {
			updatedAtField = &ModelField{
				Name:   field.Name,
				DBName: gormSchema[field.Name],
			}
		}
	}
	return &MetaModel{
		SoftDeletedField: softDeletedField,
		CreatedAtField:   createdAtField,
		UpdatedAtField:   updatedAtField,
	}
}
