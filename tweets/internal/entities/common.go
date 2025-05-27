package entities

type (
	PagedResponse struct {
		Total   int64       `json:"total"`
		HasMore bool        `json:"has_more"`
		Data    interface{} `json:"data"`
	}
)
