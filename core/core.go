package core

import (
	"reflect"
	"sync"

	"gorm.io/gorm/schema"
)

type Core struct {
}

// ExactModelName returns the exact model name of the struct
// It will return the name of the struct if the model is a pointer to a struct
func (c Core) ExactModelName(model any) string {
	typeOf := reflect.TypeOf(model)
	return typeOf.Elem().Name()
}

func (c Core) IsPointeOfStruct(model any) bool {
	// check if model is struct
	typeOf := reflect.TypeOf(model)
	if typeOf.Kind() != reflect.Ptr {
		return false
	}

	if typeOf.Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}

func (c Core) DetectModelInUse(listModels []*Model, modelName string) (*Model, bool) {
	for _, model := range listModels {
		if model.Name == modelName {
			return model, true
		}
	}
	return nil, false
}

func (c Core) ExactSchemaGorm(model any) map[string]string {
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic("failed to create schema")
	}

	m := make(map[string]string)
	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := field.Name
		m[modelName] = dbName
	}

	return m
}
