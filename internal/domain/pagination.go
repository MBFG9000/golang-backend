package domain

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"totalCount"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
}
