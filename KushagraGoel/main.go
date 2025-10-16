package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"KushagraGoel/internal/db/repository"
	"KushagraGoel/internal/db/usecase"
)

func uploadHandler(w http.ResponseWriter, r *http.Request, carRunUseCase *usecase.CarRunUseCase) {
	ctx := r.Context()

	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check extension
	ext := filepath.Ext(header.Filename)
	if ext != ".mcap" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Unsupported file type: %s\nPlease upload a .mcap file.", ext)
		return
	}

	// Use the use case to create a new CarRun record
	carRun, err := carRunUseCase.Create(ctx, "", "", "", "", "", "", "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create CarRun: %v", err), http.StatusInternalServerError)
		return
	}

	// Success response
	fmt.Fprintf(w,
		"MCAP uploaded successfully!\nFile name: %s\nFile size: %d bytes\nCarRun ID: %s\nDate Uploaded: %s",
		header.Filename, header.Size, carRun.ID.Hex(), carRun.Date_uploaded,
	)
}

func main() {
	godotenv.Load()
	ctx := context.Background()

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	db := client.Database("car_run_app")

	// Initialize repository and use case
	carRunRepo := repository.NewMongoCarRepository(db)
	carRunUseCase := usecase.NewCarRunUseCase(carRunRepo)

	// Chi router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Route for file uploads
	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadHandler(w, r, carRunUseCase)
	})

	fmt.Println("🚀 Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
