package main

// ==================================================
// REPOSITORY INTERFACE
// ==================================================
type UserRepository interface {
	List() []User
	FindByID(id int) (User, bool)
	Create(name string, email string) User
	UpdateByID(id int, name string, email string) (User, bool)
	DeleteByID(id int) bool
}
