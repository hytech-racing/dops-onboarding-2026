package usecase

import (
	"context"
	"time"

	"main/internal/db/repository"
	"main/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarRunUseCase struct {
	carRunRepo repository.CarRunRepository
}

func NewCarRunUseCase(repo repository.CarRunRepository) *CarRunUseCase {
	return &CarRunUseCase{
		carRunRepo: repo,
	}
}

// CreateCarRunUseCase creates a minimal CarRun with only ID and DateUploaded
func (uc *CarRunUseCase) CreateCarRunUseCase(ctx context.Context) (*models.CarRun, error) {
	
	carRun := &models.CarRun{
		ID:           primitive.NewObjectID(),
		DateUploaded: time.Now().UTC(),
		// Other fields left as zero values (empty strings)

	}

	err := uc.carRunRepo.Create(ctx, carRun)
	if err != nil {
		return nil, err
	}

	return carRun, nil
}
