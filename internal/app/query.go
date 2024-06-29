package app

import "github.com/ItsNotGoodName/ipcmanview/pkg/pagination"

type PageQuery struct {
	Page    int `query:"page" minimum:"1" default:"1"`
	PerPage int `query:"perPage" minimum:"1" maximum:"100" default:"20"`
}

func (q PageQuery) GetPage() pagination.Page {
	return pagination.Page{
		Page:    q.Page,
		PerPage: q.PerPage,
	}
}

type Order string

func (o Order) Ascending() bool {
	return o == "ascending"
}

func (o Order) Descending() bool {
	return o == "descending"
}

type OrderQuery struct {
	Order Order `query:"order" enum:"ascending,descending"`
}
