package dom

type Error struct {
	Message string
	Code    string
}

type ErrorResponse struct {
	Error   []Error `json:"error,omitempty"`
	Success bool    `json:"success"`
}
