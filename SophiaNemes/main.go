package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".mcap") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType) 

		w.Write([]byte(`{
    "error": "Unsupported file type",
    "message": "The uploaded file type is not allowed. Please upload a .mcap file."
}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse := fmt.Sprintf(`{
    "message": "MCAP uploaded successfully",
    "file": {
        "name": "%s",
        "size": %d
    }
}`, filename, handler.Size)

	w.Write([]byte(jsonResponse))
}
