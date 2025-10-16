package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"context"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"SophiaNemes/internal/db/repository"
	"SophiaNemes/internal/db/usecase"
	"github.com/joho/godotenv"

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
	godotenv.Load()
	ctx := context.Background()
	mongoURI := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("dops")

	carRunRepo := repository.NewMongoCarRunRepository(db)
	carRunUseCase := usecase.NewCarRunUseCase(carRunRepo)
	
	r := chi.NewRouter()     // create a chi router
	r.Use(middleware.Logger) // logs requests and responses
	r.Get("/", rootHandler)  // declare route and handler
	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadHandler(w, r, carRunUseCase)
	})
	http.ListenAndServe(":3000", r) // start up server at :3000
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request, carRunUseCase *usecase.CarRunUseCase) {
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

	ctx := context.Background()
	_, err = carRunUseCase.CreateCarRunUseCase(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorResp := McapError{
			Err:     "Database error",
			Message: "Failed to create car run record",
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