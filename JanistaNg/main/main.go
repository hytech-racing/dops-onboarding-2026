package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
func main() {
 r := chi.NewRouter() // create a chi router 
 r.Use(middleware.Logger) // logs requests and responses
 r.Get("/", rootHandler) // declare route and handler
 r.Post("/upload", uploadHandler)
 http.ListenAndServe(":3000", r) // start up server at :3000
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
 w.Write([]byte("Hello World"))
}
func uploadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

	if !strings.HasSuffix(header.Filename, ".mcap"){
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"error": "Unsupported file type", 
			"message": "The uploaded file type is not allowed. Please upload an .mcap file."
		}`))
		return // don't continue parsing request after error'ed
	}

	fileSize, err := io.ReadAll(file)
	if err != nil{
		http.Error(w, "Error reading file", http.StatusInternalServerError)
        return
	}
	if len(fileSize) != int(header.Size) {
		http.Error(w, "File size mismatch, file uploaded is not file received", http.StatusBadRequest)
		return
	}

	// if nothing errors --> success msg
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{
		"message": "MCAP uploaded successfully",
        "file": {
            "name": "%s",
            "size": %d
        }
	}`, header.Filename, header.Size)))
}