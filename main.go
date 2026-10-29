package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// ==================================================
// 1. MODEL DAN INPUT
// ==================================================
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

// ==================================================
// 2. BUSINESS ERRORS
// ==================================================
var ErrUserNotFound = errors.New("user tidak ditemukan")
var ErrNameRequired = errors.New("name wajib diisi")
var ErrEmailRequired = errors.New("email wajib diisi")

// ==================================================
// 3. REPOSITORY INTERFACE
// ==================================================
type UserRepository interface {
	List() []User
	FindByID(id int) (User, bool)
	Create(name string, email string) User
}

// ==================================================
// 4. MEMORY REPOSITORY
// ==================================================
type MemoryUserRepository struct {
	mutex  sync.RWMutex
	users  []User
	nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: []User{
			{
				ID:    1,
				Name:  "Imam",
				Email: "imam@example.com",
			},
			{
				ID:    2,
				Name:  "Syafii",
				Email: "syafii@example.com",
			},
		},
		nextID: 3,
	}
}

func (repository *MemoryUserRepository) List() []User {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	users := make([]User, len(repository.users))
	copy(users, repository.users)

	return users
}

func (repository *MemoryUserRepository) FindByID(id int) (User, bool) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	for _, user := range repository.users {
		if user.ID == id {
			return user, true
		}
	}

	return User{}, false
}

func (repository *MemoryUserRepository) Create(name string, email string) User {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	user := User{
		ID:    repository.nextID,
		Name:  name,
		Email: email,
	}
	repository.nextID++
	repository.users = append(repository.users, user)

	return user
}

// ==================================================
// 5. SERVICE
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

func (service UserService) List() []User {
	return service.repository.List()
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

// ==================================================
// 6. HTTP HANDLER
// Handler sekarang menggunakan service, bukan repository.
// ==================================================
type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler UserHandler) ListHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	users := handler.service.List()

	writeJSON(w, http.StatusOK, users)
}

func (handler UserHandler) GetByIDHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	idText := r.PathValue("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "User ID tidak valid",
		})
		return
	}

	user, err := handler.service.GetByID(id)
	if errors.Is(err, ErrUserNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Data tidak ditemukan",
		})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (handler UserHandler) CreateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Data JSON tidak valid",
		})
		return
	}

	user, err := handler.service.Create(CreateUserInput{
		Name:  input.Name,
		Email: input.Email,
	})

	if errors.Is(err, ErrNameRequired) ||
		errors.Is(err, ErrEmailRequired) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

// ==================================================
// 7. JSON HELPER
// ==================================================

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode data: %v", err)
	}
}

// ==================================================
// 8. MAIN
// Membuat dan menghubungkan semua dependency.
// ==================================================
func main() {
	repository := NewMemoryUserRepository()
	service := NewUserService(repository)
	userHandler := NewUserHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.ListHandler)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByIDHandler)
	mux.HandleFunc("POST /users", userHandler.CreateHandler)

	log.Println("server berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
