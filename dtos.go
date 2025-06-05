package crud_generator

import (
	"net/http"
	"strconv"

	"github.com/duytacong24895/go-curd-generator/core"
)

type GetListQueryParams struct {
	Page     int          `json:"page" form:"page"`
	PageSize int          `json:"page_size" form:"page_size"`
	Filter   core.IFilter `json:"filter" form:"filter"`
	OrderBy  string       `json:"order_by" form:"order_by"`
}

func (g *GetListQueryParams) Bind(r *http.Request) error {
	g.OrderBy = r.URL.Query().Get("order_by")

	rPage := r.URL.Query().Get("page")
	if rPage != "" {
		var err error
		g.Page, err = strconv.Atoi(rPage)
		if err != nil {
			return err
		}
	} else {
		g.Page = 0
	}

	rPageSize := r.URL.Query().Get("page_size")
	if rPageSize != "" {
		var err error
		g.PageSize, err = strconv.Atoi(rPageSize)
		if err != nil {
			return err
		}
	} else {
		g.PageSize = 0
	}

	strFilterInput := r.URL.Query().Get("filter")
	g.Filter = core.NewFilter()
	if err := g.Filter.Load(strFilterInput); err != nil {
		return err
	}

	return nil
}

type GetListResponse struct {
	Data       []*map[string]any `json:"data"`
	TotalCount int               `json:"total_count"`
}
