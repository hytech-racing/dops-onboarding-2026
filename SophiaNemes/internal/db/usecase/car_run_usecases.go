package usecase

import(
	"context"
	"fmt"
	"SophiaNemes/internal/db/repository"
	"SophiaNemes/models"
)

type CarRunUseCase struct {
	carRunRepo repository.CarRunRepository
}

func NewCarRunUseCase(carRunRepo repository.CarRunRepository) *CarRunUseCase {
	return &CarRunUseCase{
		carRunRepo: carRunRepo, }
}

func (uc *CarRunUseCase) CreateCarRunUseCase(ctx context.Context) (*models.CarRun, error) {
	carRun, err := models.NewCarRun()
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = uc.carRunRepo.Create(ctx, carRun)
	if err != nil {
		return nil, fmt.Errorf("failed to create car run: %w", err)
	}

	return carRun, nil
}