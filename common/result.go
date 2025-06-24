package common

type Result[T any] struct {
	Data []T  `json:"data"`
	Error string `json:"error,omitempty"`
	Success bool `json:"success"`
}