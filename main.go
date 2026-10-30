package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// ==================================================
// HTTP HANDLER
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

func (handler UserHandler) UpdateHandler(
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

	var request UpdateUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Data JSON tidak valid",
		})
		return
	}

	user, err := handler.service.UpdateByID(UpdateUserInput{
		ID:    id,
		Name:  request.Name,
		Email: request.Email,
	})

	switch {
	case errors.Is(err, ErrNameRequired),
		errors.Is(err, ErrEmailRequired):
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return

	case errors.Is(err, ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return

	case err != nil:
		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "internal server error",
			},
		)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (handler UserHandler) DeleteHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	idText := r.PathValue("id")

	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "User ID tidak valid",
		})
		return
	}

	err = handler.service.DeleteByID(id)

	if errors.Is(err, ErrUserNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": ErrUserNotFound.Error(),
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

	w.WriteHeader(http.StatusNoContent)
}

// ==================================================
// JSON HELPER
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
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateHandler)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteHandler)

	log.Println("server berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
