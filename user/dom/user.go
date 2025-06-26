package dom

import (
	"time"
)

type User struct {
	Id         string     `json:"id"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	LastAccess *time.Time `json:"lastAccess"`
}
