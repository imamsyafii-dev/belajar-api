package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
