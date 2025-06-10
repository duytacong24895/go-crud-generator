package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	constants "github.com/duytacong24895/go-crud-generator/const"
	"github.com/duytacong24895/go-crud-generator/core"
	"github.com/duytacong24895/go-crud-generator/dtos"
	"github.com/duytacong24895/go-crud-generator/services"
)

type Handler struct {
	Service      services.IService
	ListModels   []*core.Model
	DTOGetDetail func(w http.ResponseWriter, r *http.Request, ref any) any
	DTOGetList   func(w http.ResponseWriter, r *http.Request, ref any, total, page, pageSize uint) any
	DTOError     func(w http.ResponseWriter, r *http.Request, err error, errMsg string) any
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	var inputData = new(dtos.GetListQueryParams)
	if err := inputData.Bind(r); err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}

	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	resData, total, err := h.Service.GetList(model, inputData)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}

	h.ResponseGetList(w, r, resData, uint(total),
		uint(inputData.Page), uint(inputData.PageSize))
}

func (h *Handler) GetListById(w http.ResponseWriter, r *http.Request) {

	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	id := chi.URLParam(r, "id")
	res, err := h.Service.GetByID(model, id)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	h.ResponseDetail(w, r, res)
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.Service.Create(model, &inputData)
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
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.Service.Update(model, &inputData, id)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	h.ResponseDetail(w, r, res)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	model, ok := r.Context().Value(constants.ModelKey).(*core.Model)
	if !ok {
		h.ResponseError(w, r, nil, "Model not found in context")
		return
	}
	id := chi.URLParam(r, "id")
	err := h.Service.Delete(model, id)
	if err != nil {
		h.ResponseError(w, r, err, err.Error())
		return
	}
	h.ResponseDetail(w, r, nil)
}

func (h *Handler) ResponseError(w http.ResponseWriter, r *http.Request,
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

func (h *Handler) ResponseDetail(w http.ResponseWriter, r *http.Request,
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

func (h *Handler) ResponseGetList(w http.ResponseWriter, r *http.Request,
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
