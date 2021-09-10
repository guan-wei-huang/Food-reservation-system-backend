package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	GetUser(ctx context.Context, name string) (*User, error)
	CreateUser(ctx context.Context, name, password string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(dsn string) (Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &repository{db}, nil
}

func (r *repository) Close() {
	r.db.Close()
}

func (r *repository) GetUser(ctx context.Context, name string) (*User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, password
		FROM user
		WHERE name = $1`,
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) CreateUser(ctx context.Context, name, password string) (*User, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user (name, password)
		VALUES ($1, $2)`,
		name,
		password,
	)
	if err != nil {
		return nil, err
	}

	uid, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       int(uid),
		Name:     name,
		Password: password,
	}, nil
}
