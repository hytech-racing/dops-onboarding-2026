package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// JSON response structs
type McapSuccess struct {
	Message string   `json:"message"`
	File    McapFile `json:"file"`
}

type McapFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type McapError struct {
	Err     string `json:"error"`
	Message string `json:"message"`
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".mcap") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)

		errorResp := McapError{
			Err:     "Unsupported file type",
			Message: "The uploaded file type is not allowed. Please upload a .mcap file.",
		}

		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fileInfo := McapFile{
		Name: filename,
		Size: header.Size,
	}

	successResp := McapSuccess{
		Message: "MCAP uploaded successfully",
		File:    fileInfo,
	}

	json.NewEncoder(w).Encode(successResp)
}