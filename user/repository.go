package user

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	GetUser(ctx context.Context, name string) (*User, error)
	CheckUserExist(ctx context.Context, name string) (bool, error)
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
		FROM users
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

func (r *repository) CheckUserExist(ctx context.Context, name string) (bool, error) {
	row, err := r.db.QueryContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE name = $1
		)`,
		name,
	)
	if err != nil {
		return false, err
	}
	defer row.Close()

	var exists bool
	for row.Next() {
		if err := row.Scan(&exists); err != nil {
			return false, err
		}
	}

	if err = row.Err(); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) CreateUser(ctx context.Context, name, password string) (*User, error) {
	result := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (name, password)
		VALUES ($1, $2)
		RETURNING id`,
		name,
		password,
	)

	var uid int
	if err := result.Scan(&uid); err != nil {
		return nil, err
	}
	return &User{
		ID:       uid,
		Name:     name,
		Password: password,
	}, nil
}
