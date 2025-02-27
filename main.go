package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Create a new router
	mux := http.NewServeMux()
	
	// API routes
	mux.HandleFunc("/api/current-time", GetCurrentTimeHandler)
	mux.HandleFunc("/api/convert-time", ConvertTimeHandler)
	
	// Determine the frontend path
	frontendPath := "../frontend/build"
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		// Try alternative path if the build directory doesn't exist
		frontendPath = "./frontend/build"
		if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
			log.Printf("Warning: Frontend build directory not found at %s or %s", "../frontend/build", "./frontend/build")
		}
	}
	
	// Create a FileServer handler for static files
	fs := http.FileServer(http.Dir(frontendPath))
	
	// Serve static files from the React build directory
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the file exists in the build directory
		path := filepath.Join(frontendPath, r.URL.Path)
		_, err := os.Stat(path)
		
		// If the file doesn't exist, serve index.html (for SPA routing)
		if os.IsNotExist(err) && r.URL.Path != "/" {
			http.ServeFile(w, r, filepath.Join(frontendPath, "index.html"))
			return
		}
		
		// Otherwise, serve the requested file
		fs.ServeHTTP(w, r)
	}))
	
	// Determine port (use environment variable PORT if available, default to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Start the server
	serverAddr := ":" + port
	fmt.Printf("Starting server on %s...\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, mux))
}
