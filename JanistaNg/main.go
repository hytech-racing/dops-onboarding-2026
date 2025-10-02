package main

import (
	"JanistaNg/internal/db/repository"
	"JanistaNg/internal/db/usecase"
	"JanistaNg/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// cp2
var CarRunUseCase *usecase.CarRunUseCase

func main() {

	//cp2
	ctx := context.Background()
	mongoURI := "mongodb://localhost:27017"
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	db := client.Database("car_runs")

	carRunRepo := repository.NewMongoCarRunRepository(db)
	carRunUseCase := usecase.NewCarRunUseCase(carRunRepo)

	newCarRun, err := carRunUseCase.CreateCarRunUseCase(ctx)

	if err != nil {
		log.Fatal("Failed to create CarRun:", err)
	}
	fmt.Println("Created CarRun:", newCarRun)

	r := chi.NewRouter()     // create a chi router
	r.Use(middleware.Logger) // logs requests and responses
	r.Get("/", rootHandler)  // declare route and handler
	r.Post("/upload", uploadHandler)
	http.ListenAndServe(":3000", r) // start up server at :3000
}

// creating json structs
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".mcap") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Header().Set("Content-Type", "application/json")
		message := McapError{
			Err:     "The uploaded file type is not allowed. Please upload a .mcap file",
			Message: "Unsupported file type",
		}
		json.NewEncoder(w).Encode(message)
		return // don't continue parsing request after error'ed
	}

	fileSize, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	if len(fileSize) != int(header.Size) {
		http.Error(w, "File size mismatch, file uploaded is not file received", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	carRun, err := CarRunUseCase.CreateCarRunUseCase(ctx)
	if err != nil {
		http.Error(w, "Failed to create CarRun: "+err.Error(), http.StatusInternalServerError)
		return
	}

	carRun.File = models.FileInfo{
		AwsBucket: "",
		FilePath:  "",
		FileName:  header.Filename,
	}

	// if nothing errors --> success msg
	w.Header().Set("Content-Type", "application/json")

	f := McapFile{
		Name: header.Filename,
		Size: int32(header.Size),
	}
	message := McapSuccess{
		Message: "MCAP uploaded successfully",
		File:    f,
	}
	json.NewEncoder(w).Encode(message)
}
