package main

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Bentuk JSON dari HTTP request.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Input untuk service.
// Tidak menggunakan JSON tag karena bukan bagian HTTP
type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserInput struct {
	ID    int
	Name  string
	Email string
}
