package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/internal/db/repository"
	"main/internal/db/usecase"
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

var carRunUsecase *usecase.CarRunUseCase

func main() {
	// ====== MongoDB Setup ======
	ctx := context.Background()
	mongoURI := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}
	defer client.Disconnect(ctx)

	db := client.Database("hytech") // or your DB name
	carRunRepo := repository.NewMongoCarRunRepository(db)
	carRunUsecase = usecase.NewCarRunUseCase(carRunRepo)

	// ====== HTTP Router ======
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/upload", uploadHandler)

	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 100 MB)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		sendJSONError(w, "Bad request", "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		sendJSONError(w, "Bad request", "No file provided in the request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file extension (simple check)
	if success, _ := isMCAP(header.Filename); !success {
		sendJSONError(w, "Unsupported file type",
			"The uploaded file type is not allowed. Please upload a .mcap file.", http.StatusUnsupportedMediaType)
		return
	}

	// ✅ CHECKPOINT 2: Create CarRun record in MongoDB
	_, err = carRunUsecase.CreateCarRunUseCase(r.Context())
	if err != nil {
		sendJSONError(w, "Internal server error", "Failed to create car run record", http.StatusInternalServerError)
		return
	}

	// Success response (no file saved, so size = 0 or omit)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadResponse{
		Message: "MCAP upload initiated – metadata recorded",
		File: FileInfo{
			Name: header.Filename,
			Size: 0, // since we're not reading/saving the file
		},
	})
}

func isMCAP(filename string) (bool, error) {
	return strings.HasSuffix(strings.ToLower(filename), ".mcap"), nil
}

func sendJSONError(w http.ResponseWriter, errorMsg, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorMsg,
		Message: message,
	})
}
