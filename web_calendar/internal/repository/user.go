package repository

import (
	"context"
	"fmt"
)

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"users_email"`
	Password string `db:"users_password"`
}

func (r *Repository) Login(ctx context.Context, email, password string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, users_email, users_password from users where users_email = $1 AND users_password = $2`, email, password)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	err = row.Scan(&u.Id, &u.Email, &u.Password)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) AddNewUser(ctx context.Context, email, password string) (u User, err error) {
	_, err = r.pool.Exec(ctx, `insert into users (users_email, users_password) values ($1, $2)`, email, password)
	if err != nil {
		err = fmt.Errorf("failed to add user: %w", err)
		return
	}
	return
}
