package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type SuccessResponse struct {
	Message string   `json:"message"`
	File    FileInfo `json:"file"`
}

func main() {
	r := chi.NewRouter()     // create a chi router
	r.Use(middleware.Logger) // logs requests and responses
	r.Get("/", rootHandler)  // declare route and handler
	r.Post("/upload", uploadHandler)
	http.ListenAndServe(":3000", r) // start up server at :3000
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Step 2: Get the file from the form (like pulling the file out of the envelope)
	file, handler, err := r.FormFile("file") // "file" is the name of the form field
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close() // Clean up when we're done

	// Step 3: Check if it's a .mcap file
	filename := handler.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".mcap") {
		// It's NOT a .mcap file, so send error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType) // 415 error

		errorResponse := ErrorResponse{
			Error:   "Unsupported file type",
			Message: "The uploaded file type is not allowed. Please upload a .mcap file.",
		}

		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Step 4: It IS a .mcap file, so send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 success

	successResponse := SuccessResponse{
		Message: "MCAP uploaded successfully",
		File: FileInfo{
			Name: filename,
			Size: handler.Size,
		},
	}

	json.NewEncoder(w).Encode(successResponse)

}
