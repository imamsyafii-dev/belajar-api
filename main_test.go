package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeUserRepository struct {
	users        []User
	nextID       int
	createdName  string
	createdEmail string
}

func (repository *FakeUserRepository) List(
	_ context.Context,
) (
	[]User,
	error,
) {
	return repository.users, nil
}

func (repository *FakeUserRepository) FindByID(
	id int,
) (User, bool, error) {
	for _, user := range repository.users {
		if user.ID == id {
			return user, true, nil
		}
	}

	return User{}, false, nil
}

func (repository *FakeUserRepository) Create(
	name string,
	email string,
) (User, error) {
	repository.createdName = name
	repository.createdEmail = email

	user := User{
		ID:    repository.nextID,
		Name:  name,
		Email: email,
	}

	repository.users = append(repository.users, user)
	repository.nextID++

	return user, nil
}

func (repository *FakeUserRepository) DeleteByID(
	id int,
) (bool, error) {
	for index, user := range repository.users {
		if user.ID == id {
			repository.users = append(
				repository.users[:index],
				repository.users[index+1:]...,
			)

			return true, nil
		}
	}

	return false, nil
}

func (repository *FakeUserRepository) UpdateByID(
	id int,
	name string,
	email string,
) (User, bool, error) {
	for index, user := range repository.users {
		if user.ID == id {
			user.Name = name
			user.Email = email

			repository.users[index] = user
			return user, true, nil
		}
	}
	return User{}, true, nil
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
		Name:  "     ",
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

func TestCreateUserHandlerSuccess(t *testing.T) {
	// Arrange: siapkan dependency.
	repository := &FakeUserRepository{
		nextID: 3,
	}

	service := NewUserService(repository)
	handler := NewUserHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", handler.CreateHandler)

	body := strings.NewReader(`{
		"name": "Citra",
		"email": "citra@example.com"
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		body,
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	response := httptest.NewRecorder()

	// Act: kirim request ke router.
	mux.ServeHTTP(response, request)

	// Assert: periksa HTTP status.
	if response.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.Code,
		)
	}

	// Decode JSON response.
	var user User

	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if user.ID != 3 {
		t.Errorf("expected ID 3, got %d", user.ID)
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
}

func TestCreateUserHandlerNameRequired(t *testing.T) {
	repository := &FakeUserRepository{
		nextID: 3,
	}

	service := NewUserService(repository)
	handler := NewUserHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", handler.CreateHandler)

	body := strings.NewReader(`{
		"name": "",
		"email": "citra@example.com"
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		body,
	)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	var result map[string]string

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if result["error"] != ErrNameRequired.Error() {
		t.Errorf(
			"expected error %q, got %q",
			ErrNameRequired.Error(),
			result["error"],
		)
	}

	if len(repository.users) != 0 {
		t.Error("repository should not create a user")
	}
}

func TestDeleteUserHandlerSuccess(t *testing.T) {
	repository := &FakeUserRepository{
		users: []User{
			{
				ID:    1,
				Name:  "Citra",
				Email: "citra@example.com",
			},
		},
	}

	service := NewUserService(repository)
	handler := NewUserHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"DELETE /users/{id}",
		handler.DeleteHandler,
	)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/users/1",
		nil,
	)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			response.Code,
		)
	}

	if len(repository.users) != 0 {
		t.Errorf(
			"expected no users, got %d",
			len(repository.users),
		)
	}
}

func TestUpdateUserHandlerSuccess(t *testing.T) {
	repository := &FakeUserRepository{
		users: []User{
			{
				ID:    1,
				Name:  "Imam",
				Email: "imam@example.com",
			},
		},
	}

	service := NewUserService(repository)
	handler := NewUserHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"PUT /users/{id}",
		handler.UpdateHandler,
	)

	body := strings.NewReader(`{
		"name": "Imam Syafii",
		"email": "imam.syafii@example.com"
	}`)

	request := httptest.NewRequest(
		http.MethodPut,
		"/users/1",
		body,
	)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	var user User

	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if user.Name != "Imam Syafii" {
		t.Errorf(
			"expected name Imam Syafii, got %s",
			user.Name,
		)
	}

	if repository.users[0].Email != "imam.syafii@example.com" {
		t.Errorf(
			"repository email was not updated: %s",
			repository.users[0].Email,
		)
	}
}
