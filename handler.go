package curd_generator

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	constants "github.com/duytacong24895/go-curd-generator/const"
	"github.com/duytacong24895/go-curd-generator/core"
)

type handler struct {
	service      IService
	listModels   []*core.Model
	core         *core.Core
	DTOGetDetail func(w http.ResponseWriter, r *http.Request, ref any) any
	DTOGetList   func(w http.ResponseWriter, r *http.Request, ref any, total, page, pageSize uint) any
	DTOError     func(w http.ResponseWriter, r *http.Request, err error, errMsg string) any
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	var inputData = new(GetListQueryParams)
	if err := inputData.Bind(r); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}

	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	resData, total, err := h.service.GetList(model, inputData)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}

	h.ResponseGetList(w, r, resData, uint(total),
		uint(inputData.Page), uint(inputData.PageSize))
}

func (h *handler) GetListById(w http.ResponseWriter, r *http.Request) {

	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := h.service.GetByID(model, id)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	h.ResponseDetail(w, r, res)
}
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	// get params
	var inputData = make(map[string]any)
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}

	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	res, err := h.service.Create(model, &inputData)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
}
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var inputData = make(map[string]any)
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}

	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := h.service.Update(model, &inputData, id)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	h.ResponseDetail(w, r, res)
}
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	id := chi.URLParam(r, "id")
	err := h.service.Delete(model, id)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	h.ResponseDetail(w, r, nil)
}

func (h *handler) ResponseError(w http.ResponseWriter, r *http.Request,
	err error, msgErr string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if h.DTOError == nil {
		http.Error(w, msgErr, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(h.DTOError(w, r, err, msgErr)); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
}

func (h *handler) ResponseDetail(w http.ResponseWriter, r *http.Request,
	ref any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if h.DTOGetDetail == nil {
		if err := json.NewEncoder(w).Encode(ref); err != nil {
			h.ResponseError(w, r, err, err.Error())
		}
		return
	}

	if err := json.NewEncoder(w).Encode(h.DTOGetDetail(w, r, ref)); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
}

func (h *handler) ResponseGetList(w http.ResponseWriter, r *http.Request,
	ref any, total, page, pageSize uint) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if h.DTOGetList == nil {
		if err := json.NewEncoder(w).Encode(ref); err != nil {
			h.ResponseError(w, r, err, err.Error())
		}
		return
	}

	if err := json.NewEncoder(w).Encode(h.DTOGetList(w, r, ref, total, page, pageSize)); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
}
