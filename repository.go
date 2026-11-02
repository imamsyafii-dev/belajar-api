package main

import "context"

// ==================================================
// REPOSITORY INTERFACE
// ==================================================
type UserRepository interface {
	List(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id int) (User, bool, error)
	Create(ctx context.Context, name string, email string) (User, error)
	UpdateByID(ctx context.Context, id int, name string, email string) (User, bool, error)
	DeleteByID(ctx context.Context, id int) (bool, error)
}
