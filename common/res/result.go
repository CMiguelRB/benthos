package common

type WResult struct {
	AffectedRows int64  `json:"rows"`
	Error string `json:"error,omitempty"`
	Success bool `json:"success"`
}

type QResult[T any] struct {
	Data []T  `json:"data"`
	Error string `json:"error,omitempty"`
	Success bool `json:"success"`
}