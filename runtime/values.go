package runtime

import (
	"sync"

	"github.com/duytacong24895/go-crud-generator/core"
)

var once sync.Once

type RegisteredModels struct {
	List []*core.Model
}

func (r *RegisteredModels) Add(model *core.Model) {
	r.List = append(r.List, model)
}

var registeredModels *RegisteredModels

func GetListModels() *RegisteredModels {
	once.Do(func() {
		// Register your custom values here
		registeredModels = &RegisteredModels{
			List: make([]*core.Model, 0),
		}
	})
	return registeredModels
}
