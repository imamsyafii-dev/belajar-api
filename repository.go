package main

// ==================================================
// REPOSITORY INTERFACE
// ==================================================
type UserRepository interface {
	List() ([]User, error)
	FindByID(id int) (User, bool, error)
	Create(name string, email string) (User, error)
	UpdateByID(id int, name string, email string) (User, bool)
	DeleteByID(id int) bool
}
