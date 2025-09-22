package main

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check extension
	ext := filepath.Ext(header.Filename)
	if ext != ".mcap" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Unsupported file type: %s\nThe uploaded file type is not allowed. Please upload a .jpg or .png file.", ext)
		return
	}

	// Success response
	fmt.Fprintf(w, "MCAP uploaded successfully\nFile name: %s\nFile size: %d bytes", header.Filename, header.Size)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
