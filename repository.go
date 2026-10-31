package main

import "context"

// ==================================================
// REPOSITORY INTERFACE
// ==================================================
type UserRepository interface {
	List(ctx context.Context) ([]User, error)
	FindByID(id int) (User, bool, error)
	Create(name string, email string) (User, error)
	UpdateByID(id int, name string, email string) (User, bool, error)
	DeleteByID(id int) (bool, error)
}
