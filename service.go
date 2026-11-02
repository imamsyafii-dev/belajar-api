package main

import (
	"context"
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

func (service UserService) List(
	ctx context.Context,
) (
	[]User,
	error,
) {
	users, err := service.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service UserService) GetByID(
	ctx context.Context,
	id int,
) (User, error) {
	user, found, err := service.repository.FindByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (service UserService) Create(
	ctx context.Context,
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

	user, err := service.repository.Create(
		ctx,
		input.Name,
		input.Email,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (service UserService) UpdateByID(
	ctx context.Context,
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

	user, found, err := service.repository.UpdateByID(
		ctx,
		input.ID,
		input.Name,
		input.Email,
	)

	if err != nil {
		return User{}, err
	}

	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (service UserService) DeleteByID(
	ctx context.Context,
	id int,
) error {
	deleted, err := service.repository.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrUserNotFound
	}

	return nil
}
