package main

import (
	"errors"
	"strings"
)

// ==================================================
// BUSINESS ERRORS
// ==================================================
var ErrUserNotFound = errors.New("user tidak ditemukan")
var ErrNameRequired = errors.New("name wajib diisi")
var ErrEmailRequired = errors.New("email wajib diisi")

// ==================================================
// SERVICE
// Bagian baru yang berada di antara handler dan repository.
// ==================================================

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (service UserService) List() (
	[]User,
	error,
) {
	users, err := service.repository.List()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service UserService) GetByID(id int) (User, error) {
	user, found := service.repository.FindByID(id)
	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (service UserService) Create(
	input CreateUserInput,
) (User, error) {
	// Validasi dipindahkan dari handler ke service.
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" {
		return User{}, ErrNameRequired
	}

	if input.Email == "" {
		return User{}, ErrEmailRequired
	}

	user := service.repository.Create(
		input.Name,
		input.Email,
	)

	return user, nil
}

func (service UserService) UpdateByID(
	input UpdateUserInput,
) (User, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" {
		return User{}, ErrNameRequired
	}

	if input.Email == "" {
		return User{}, ErrEmailRequired
	}

	user, found := service.repository.UpdateByID(
		input.ID,
		input.Name,
		input.Email,
	)
	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (service UserService) DeleteByID(id int) error {
	deleted := service.repository.DeleteByID(id)
	if !deleted {
		return ErrUserNotFound
	}

	return nil
}
