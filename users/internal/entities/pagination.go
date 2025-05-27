package entities

type (
	Pagination struct {
		Page     int64 `json:"page"`
		PageSize int64 `json:"page_size"`
	}

	PagedResponse struct {
		Total   int64       `json:"total"`
		HasMore bool        `json:"has_more"`
		Data    interface{} `json:"data"`
	}
)

func NewDefaultPagination() Pagination {
	return Pagination{
		Page:     1,
		PageSize: 25,
	}
}

func (p *Pagination) SafePageNumber() int64 {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}
