package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/db/repository"
	"main/internal/db/usecase"
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

	ctx := context.Background()
	mongoURI := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic("Error: could not create repo client")

	}
	db := client.Database("vehicle_data_db")

	carRunRepo, err := repository.NewMongoCarRepository(db)
	if err != nil {
		panic("Error: could not create repo")
	}

	carRunUseCase := usecase.NewCarRunUseCase(carRunRepo)

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		UploadMcap(ctx, carRunUseCase, w, r)
	})

	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", r)
}

func UploadMcap(ctx context.Context, carRunUseCase *usecase.CarRunUseCase, w http.ResponseWriter, r *http.Request) {
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
