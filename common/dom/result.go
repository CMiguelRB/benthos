package dom

type WResult struct {
	Rows    *int64  `json:"rows,omitempty"`
	Id      *string `json:"id,omitempty"`
	Error   string  `json:"error,omitempty"`
	Success bool    `json:"success"`
}

type QResult[T any] struct {
	Data    []T    `json:"data"`
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}
