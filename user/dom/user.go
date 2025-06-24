package dom

import (
	"time"
)

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	CreatedOn *time.Time `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	LastAccess *time.Time `json:"lastAccess"`
}