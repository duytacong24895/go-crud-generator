package crud_generator

import (
	"context"
	"net/http"

	constants "github.com/duytacong24895/go-curd-generator/const"
	"github.com/duytacong24895/go-curd-generator/core"
	"github.com/go-chi/chi/v5"
)

func VerifyModel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var inputModelname = chi.URLParam(r, "modelName")
		var model, ok = core.Core{}.DetectModelInUse(listModels, inputModelname)
		if !ok {
			http.Error(w, "Model not found", http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), constants.ModelKey, model)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
