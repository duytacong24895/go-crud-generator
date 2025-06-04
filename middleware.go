package crud_generator

import (
	"context"
	"net/http"

	constants "github.com/Vietnam-Silicon/template-api-go/plugins/crud_generator/const"
	"github.com/Vietnam-Silicon/template-api-go/plugins/crud_generator/core"
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
