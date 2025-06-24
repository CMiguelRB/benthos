package infra

import (
	"context"
	"log/slog"
	"time"

	sec "benthos/common/sec"
	"benthos/db"
	"benthos/user/dom"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	getUsersQuery    string
	getUserByIdQuery string
	createUserQuery  string
	updateUserQuery  string
	deleteUserQuery  string
}

func NewRepo() *Repo {
	return &Repo{
		getUsersQuery:    "SELECT * FROM users ORDER BY created_on ASC",
		getUserByIdQuery: "SELECT * FROM users WHERE id = $1",
		createUserQuery:  "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		updateUserQuery:  "UPDATE users SET username = $1, password = $2, updated_on = $3 WHERE id = $4",
		deleteUserQuery:  "DELETE FROM users WHERE id = $1",
	}
}

func (r *Repo) GetUsers(ctx context.Context) (users []dom.User, error error) {

	rows, err := db.Pool.Query(ctx, r.getUsersQuery)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[dom.User])

	for i := 0; i < len(users); i++ {
		decryptedPassword, err := sec.Decrypt(users[i].Password)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		users[i].Password = decryptedPassword
	}

	return users, err
}

func (r *Repo) GetUserById(ctx context.Context, id string) (user []dom.User, error error) {

	rows, err := db.Pool.Query(ctx, r.getUserByIdQuery, id)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	user, err = pgx.CollectRows(rows, pgx.RowToStructByName[dom.User])

	for i := 0; i < len(user); i++ {
		decryptedPassword, err := sec.Decrypt(user[i].Password)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		user[i].Password = decryptedPassword
	}

	return user, err
}

func (r *Repo) CreateUser(ctx context.Context, user dom.User) (string, error) {

	password, err := sec.Encrypt(user.Password)

	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	var id string
	err = db.Pool.QueryRow(ctx, r.createUserQuery, user.Username, password).Scan(&id)

	if err != nil {
		slog.Error(err.Error())
	}

	return id, err
}

func (r *Repo) UpdateUser(ctx context.Context, id string, user dom.User) (int64, error) {

	password, err := sec.Encrypt(user.Password)

	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	datetime := time.Now()

	res, err := db.Pool.Exec(ctx, r.updateUserQuery, user.Username, password, datetime, id)

	if err != nil {
		slog.Error(err.Error())
	}

	return res.RowsAffected(), err
}

func (r *Repo) DeleteUser(ctx context.Context, id string) (int64, error) {

	res, err := db.Pool.Exec(ctx, r.deleteUserQuery, id)

	if err != nil {
		slog.Error(err.Error())
	}

	return res.RowsAffected(), err
}
