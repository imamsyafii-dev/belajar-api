package main

import (
	"errors"
	"testing"
)

type FakeUserRepository struct {
	users        []User
	nextID       int
	createdName  string
	createdEmail string
}

func (repository *FakeUserRepository) List() []User {
	return repository.users
}

func (repository *FakeUserRepository) FindByID(
	id int,
) (User, bool) {
	for _, user := range repository.users {
		if user.ID == id {
			return user, true
		}
	}

	return User{}, false
}

func (repository *FakeUserRepository) Create(
	name string,
	email string,
) User {
	repository.createdName = name
	repository.createdEmail = email

	user := User{
		ID:    repository.nextID,
		Name:  name,
		Email: email,
	}

	repository.users = append(repository.users, user)
	repository.nextID++

	return user
}

func TestUserServiceCreateSuccess(t *testing.T) {
	repository := &FakeUserRepository{
		nextID: 10,
	}

	service := NewUserService(repository)
	user, err := service.Create(CreateUserInput{
		Name:  " Citra ",
		Email: " citra@example.com ",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 10 {
		t.Errorf("expected ID 10, got %d", user.ID)
	}

	if user.Name != "Citra" {
		t.Errorf(
			"expected name Citra, got %s",
			user.Name,
		)
	}

	if user.Email != "citra@example.com" {
		t.Errorf(
			"expected email citra@example.com, got %s",
			user.Email,
		)
	}

	if repository.createdName != "Citra" {
		t.Errorf(
			"expected repository name Citra, got %s",
			repository.createdName,
		)
	}
}

func TestUserServiceCreateNameRequired(t *testing.T) {
	repository := &FakeUserRepository{
		nextID: 1,
	}

	service := NewUserService(repository)

	_, err := service.Create(CreateUserInput{
		Name: "     ",
		Email: "citra@example.com",
	})

	if !errors.Is(err, ErrNameRequired) {
		t.Fatalf(
			"expected ErrNameRequired, got %v",
			err,
		)
	}

	if len(repository.users) != 0 {
		t.Error("repository should not create a user")
	}
}

func TestUserServiceGetByIDNotFound(t *testing.T) {
	repository := &FakeUserRepository{
		users: []User{},
	}

	service := NewUserService(repository)

	_, err := service.GetByID(99)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}
