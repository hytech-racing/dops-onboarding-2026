package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UploadResponse struct {
	Message string   `json:"message"`
	File    FileInfo `json:"file"`
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//The logger will log information about
	// incoming requests like the request method,
	// path, and the response status.

	os.MkdirAll("./uploads", os.ModePerm)

	r.Post("/upload", uploadHandler)

	// After that, you need to set up a route to
	// the root path that listens for GET
	// requests and returns an OK back to the client:

	//starts server
	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		sendJSONError(w, "Bad request", "Unable to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		sendJSONError(w, "Bad request", "No file provided in the request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if success, err := isMCAP(header.Filename); err != nil || success != true {
		sendJSONError(w, "Unsupported file type",
			"The uploaded file type is not allowed. Please upload a .mcap file.", http.StatusUnsupportedMediaType)
		return
	}

	// Create destination file
	dstPath := filepath.Join("./uploads", header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		sendJSONError(w, "Internal server error", "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content and get size
	size, err := io.Copy(dst, file)
	if err != nil {
		sendJSONError(w, "Internal server error", "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadResponse{
		Message: "MCAP uploaded successfully",
		File: FileInfo{
			Name: header.Filename,
			Size: size,
		},
	})
}

// IsMCAPFile checks if a file is an MCAP file by examining its magic bytes.
func isMCAP(filename string) (bool, error) {
	if strings.HasSuffix(strings.ToLower(filename), ".mcap") {
		return true, nil
	}
	// file, err := os.Open(filename)
	// if err != nil {
	// 	return false, err
	// }
	// defer file.Close()

	// reader, err := mcap.NewReader(file)
	// if err != nil {
	// 	return false, nil // Not an MCAP file or corrupted
	// }
	// defer reader.Close()

	// // If we get here without error, it's a valid MCAP file
	return false, nil
}

// Helper function for error responses
func sendJSONError(w http.ResponseWriter, errorMsg, message string, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorMsg,
		Message: message,
	})

}
