package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create a new router
	mux := http.NewServeMux()
	
	// API routes
	mux.HandleFunc("/api/current-time", GetCurrentTimeHandler)
	mux.HandleFunc("/api/convert-time", ConvertTimeHandler)
	
	// Serve static files from the frontend build directory
	fs := http.FileSystem(http.Dir("../frontend/build"))
	mux.Handle("/", http.StripPrefix("/", http.FileServer(fs)))
	
	// Wrap the mux with CORS middleware
	corsHandler := enableCORS(mux)
	
	// Start the server
	fmt.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
