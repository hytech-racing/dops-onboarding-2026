package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internals/db/repository"
	"main/internals/db/usecase"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
		return
	}
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

	if strings.Contains(header.Filename, ".mcap") || header.Header.Get("Content-Type") == "application/mcap" {
		f := McapFile{
			Name: header.Filename,
			Size: int32(header.Size),
		}
		message := McapSuccess{
			Message: "MCAP uploaded successfully",
			File:    f,
		}

		ctx := context.Background()
		mongoURI := os.Getenv("MONGODB_URI")
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			http.Error(w, "Error: could not create mongo client "+err.Error(), http.StatusInternalServerError)

		}
		db := client.Database("vehicle_data_db")

		carRunRepo, err := repository.NewMongoCarRepository(db)
		if err != nil {
			http.Error(w, "Error: could not create repository "+err.Error(), http.StatusInternalServerError)
		}

		carRunUseCase := usecase.NewCarRunUseCase(carRunRepo)

		newCarRun, err := carRunUseCase.CreateCarRun(ctx)
		if err != nil {
			http.Error(w, "Error: could not insert into mongoDB "+err.Error(), http.StatusInternalServerError)
		}

		fmt.Println(*newCarRun)

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
