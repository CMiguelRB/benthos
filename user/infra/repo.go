package infra

import (
	"context"
	"log/slog"
	"time"

	"benthos_go/db"
	"benthos_go/user/dom"
	"benthos_go/common"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GetUsers(ctx context.Context) (users []dom.User, error error) {

	rows, err := db.Pool.Query(ctx, "SELECT * from users;")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[dom.User])

	for  i := 0; i<len(users); i++{
		decryptedPassword, err := common.Decrypt(users[i].Password)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		users[i].Password = decryptedPassword
	}

	return users, err
}

func (r *Repo) GetUserById(ctx context.Context, id string) (user []dom.User, error error) {

	rows, err := db.Pool.Query(ctx, "SELECT * from users where id = '"+id+"';")

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	user, err = pgx.CollectRows(rows, pgx.RowToStructByName[dom.User])

	for  i := 0; i<len(user); i++{
		decryptedPassword, err := common.Decrypt(user[i].Password)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		user[i].Password = decryptedPassword
	}

	return user, err
}

func (r *Repo) CreateUser(ctx context.Context, user dom.User) (int64, error) {
	query := `INSERT INTO users (username, password, "createdOn") VALUES (@username, @password, @createdOn)`

	password, err := common.Encrypt(user.Password)

	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	args := pgx.NamedArgs{
		"username": user.Username,
		"password": password,
		"createdOn": time.Now(),
	}

	res, err := db.Pool.Exec(ctx, query, args)
	if err != nil {
		slog.Error(err.Error())
	}

	return res.RowsAffected(), err
}

func (r *Repo) UpdateUser(ctx context.Context, id string, user dom.User) (int64, error) {
	query := `UPDATE users SET username = @username, password = @password, "updatedOn" = @updatedOn WHERE id = @id`

	password, err := common.Encrypt(user.Password)

	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	args := pgx.NamedArgs{
		"username": user.Username,
		"password": password,
		"updatedOn": time.Now(),
		"id": id,
	}

	res, err := db.Pool.Exec(ctx, query, args)
	if err != nil {
		slog.Error(err.Error())
	}

	return res.RowsAffected(), err
}

func (r *Repo) DeleteUser(ctx context.Context, id string) (int64, error){
	query := `DELETE FROM users WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	res, err := db.Pool.Exec(ctx, query, args)

	if err != nil {
		slog.Error(err.Error())
	}

	return res.RowsAffected(), err

}
