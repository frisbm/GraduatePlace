package pagination

type Pagination[T any] struct {
	Data   []*T  `json:"data"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	Count  int32 `json:"count"`
}
