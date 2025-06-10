package crud_generator

import (
	"fmt"
	"net/http"

	"github.com/duytacong24895/go-crud-generator/core"
	"github.com/duytacong24895/go-crud-generator/handler"
	"github.com/duytacong24895/go-crud-generator/middlewares"
	"github.com/duytacong24895/go-crud-generator/repositories"
	"github.com/duytacong24895/go-crud-generator/runtime"
	"github.com/duytacong24895/go-crud-generator/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ICRUDGenerator interface {
	Run()
	RegisterModel(ref any) ICRUDGenerator
	RegisterMiddleware(middleware func(next http.Handler) http.Handler) ICRUDGenerator
	RegisterDTOForGetDetail(func(w http.ResponseWriter, r *http.Request, ref any) any) ICRUDGenerator
	RegisterDTOForGetList(func(w http.ResponseWriter, r *http.Request, ref any, total, page, pageSize uint) any) ICRUDGenerator
	RegisterDTOForError(func(w http.ResponseWriter, r *http.Request, err error, errMsg string) any) ICRUDGenerator
}

type crudGenerator struct {
	router      *chi.Mux
	handler     *handler.Handler
	core        *core.Core
	middlewares []func(next http.Handler) http.Handler
}

func NewCRUDGenerator(router *chi.Mux, db *gorm.DB) ICRUDGenerator {
	core := &core.Core{}
	return &crudGenerator{
		router: router,
		core:   core,
		handler: &handler.Handler{
			Service:    services.NewService(repositories.NewRepository(db)),
			ListModels: runtime.GetListModels().List,
		},
	}
}

func (c *crudGenerator) RegisterModel(model any) ICRUDGenerator {
	if c.core.IsPointeOfStruct(model) {
		runtime.GetListModels().Add(core.NewModel(model))
	} else {
		panic("Model must be a pointer to a struct")
	}
	return c
}

func (c *crudGenerator) RegisterMiddleware(middleware func(next http.Handler) http.Handler) ICRUDGenerator {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

func (c *crudGenerator) RegisterDTOForGetDetail(returndto func(w http.ResponseWriter, r *http.Request, ref any) any) ICRUDGenerator {
	c.handler.DTOGetDetail = returndto
	return c
}

func (c *crudGenerator) RegisterDTOForGetList(returndto func(w http.ResponseWriter, r *http.Request, ref any, total, page, pageSize uint) any) ICRUDGenerator {
	c.handler.DTOGetList = returndto
	return c
}

func (c *crudGenerator) RegisterDTOForError(returndto func(w http.ResponseWriter, r *http.Request, err error, errMsg string) any) ICRUDGenerator {
	c.handler.DTOError = returndto
	return c
}

func (c *crudGenerator) Run() {
	c.router.Route("/crud", func(r chi.Router) {
		r = r.With(middlewares.VerifyModel)
		for _, middleware := range c.middlewares {
			r = r.With(middleware)
		}
		r.Get("/{modelName}", c.handler.GetList)
		r.Get("/{modelName}/{id}", c.handler.GetListById)
		r.Post("/{modelName}", c.handler.Create)
		r.Put("/{modelName}/{id}", c.handler.Update)
		r.Delete("/{modelName}/{id}", c.handler.Delete)
	})

	listModelNames := make([]string, len(runtime.GetListModels().List))
	for i, model := range runtime.GetListModels().List {
		listModelNames[i] = model.Name
	}
	fmt.Println("CRUD generator initialized")
	fmt.Println("List models:", listModelNames)
	fmt.Println("CRUD generator routes registered")
}
