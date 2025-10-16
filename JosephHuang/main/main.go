package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type McapError struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

type McapSuccess struct {
	Message string   `json:"message"`
	File    McapFile `json:"file"`
}

type McapFile struct {
	Name string `json:"name"`
	Size int32  `json:"size"`
}

var ByteOffset int = 20

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server running on :3000	"))
	})

	r.Post("/upload", UploadMcap)

	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", r)
}

func UploadMcap(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << ByteOffset)
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}

	_, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadGateway)
		return
	}

	fmt.Println(header.Header)
	fmt.Println(header.Filename)
	fmt.Println(header.Size)

	if strings.Contains(header.Filename, ".mcap") || header.Header.Get("Content-Type") == "application/mcap" {
		f := McapFile{
			Name: header.Filename,
			Size: int32(header.Size),
		}
		message := McapSuccess{
			Message: "MCAP uploaded successfully",
			File:    f,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
		return

	} else {
		errorMessage := McapError{
			Err:     "Unsupported file type",
			Message: "The uploaded file type is not allowed. Please upload a .mcap file",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

}
