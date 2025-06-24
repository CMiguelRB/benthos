package dom

import (
	"time"
	"database/sql"
)

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	CreatedOn time.Time `json:"createdOn"`
	UpdatedOn sql.NullTime `json:"updatedOn.Time"`
	LastAccess sql.NullTime `json:"lastAccess"`
}