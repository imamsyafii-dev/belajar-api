package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(
	pool *pgxpool.Pool,
) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

func (repository *PostgresUserRepository) List(
	ctx context.Context,
) (
	[]User,
	error,
) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT id, name, email
		FROM users
		ORDER BY id
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (repository *PostgresUserRepository) Create(
	ctx context.Context,
	name string,
	email string,
) (
	User,
	error,
) {
	var user User

	err := repository.pool.QueryRow(
		ctx,
		`
			INSERT INTO users (name, email)
			VALUES ($1, $2)
			RETURNING id, name, email
		`,
		name,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return User{}, fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	return user, nil
}

func (repository *PostgresUserRepository) FindByID(
	ctx context.Context,
	id int,
) (
	User,
	bool,
	error,
) {
	var user User

	err := repository.pool.QueryRow(
		ctx,
		`
			SELECT id, name, email
			FROM users
			WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err == pgx.ErrNoRows {
		return User{}, false, nil
	}

	if err != nil {
		return User{}, false, fmt.Errorf(
			"query user by id: %w", err,
		)
	}

	return user, true, nil
}

func (repository *PostgresUserRepository) UpdateByID(
	ctx context.Context,
	id int,
	name string,
	email string,
) (
	User,
	bool,
	error,
) {
	var user User

	err := repository.pool.QueryRow(
		ctx,
		`
			UPDATE users
			SET name = $1, email = $2
			WHERE id = $3
			RETURNING id, name, email
		`,
		name,
		email,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err == pgx.ErrNoRows {
		return User{}, false, nil
	}

	if err != nil {
		return User{}, false, fmt.Errorf(
			"update user: %w",
			err,
		)
	}

	return user, true, nil
}

func (repository *PostgresUserRepository) DeleteByID(
	ctx context.Context,
	id int,
) (
	bool,
	error,
) {
	commandTag, err := repository.pool.Exec(
		ctx,
		`
			DELETE FROM users
			WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return false, fmt.Errorf(
			"delete user: %w",
			err,
		)
	}

	return commandTag.RowsAffected() == 1, nil
}
