package main

import (
	"context"
	"fmt"

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
