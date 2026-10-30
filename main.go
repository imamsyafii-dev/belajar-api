package main

import (
	"log"
	"net/http"
)

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
